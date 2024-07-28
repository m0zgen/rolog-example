package rotatinglogger

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type LogMessage struct {
	Level   logrus.Level
	Message string
	Fields  logrus.Fields
}

type RotatingLogger struct {
	Logger         *logrus.Logger
	checkInterval  time.Duration
	maxSize        int
	zippedArchive  bool
	logDir         string
	archivePattern string
	logChannel     chan LogMessage
	logLevel       logrus.Level
	maxBackups     int
	staticFilePath string
	multiWriter    *bufio.Writer
	consoleOut     bool
}

func NewRotatingLogger(logDir, staticFilename, archivePattern string, zippedArchive bool, maxSize, maxBackups int, checkInterval time.Duration, bufferSize int, logLevel logrus.Level, consoleOut bool) *RotatingLogger {
	logger := logrus.New()

	logger.SetLevel(logLevel) // Set the logging level
	logger.Formatter = &JournalctlFormatter{}

	// Ensure log directory exists
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		logger.Fatalf("Failed to create log directory: %v", err)
	}

	staticFilePath := filepath.Join(logDir, staticFilename)

	// Set the logger output to the file
	file, err := os.OpenFile(staticFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	//logger.SetOutput(file)

	//multiWriter := io.Writer(file)
	multiWriter := bufio.NewWriter(file)
	if consoleOut {
		fWriter := io.Writer(file)
		mWriter := io.MultiWriter(fWriter, os.Stdout)
		logger.SetOutput(mWriter)
	} else {
		logger.SetOutput(multiWriter)
	}

	rotLogger := &RotatingLogger{
		Logger:         logger,
		checkInterval:  checkInterval,
		maxSize:        maxSize,
		zippedArchive:  zippedArchive,
		logDir:         logDir,
		archivePattern: archivePattern,
		logChannel:     make(chan LogMessage, bufferSize),
		logLevel:       logLevel,
		maxBackups:     maxBackups,
		staticFilePath: staticFilePath,
		multiWriter:    multiWriter,
		consoleOut:     consoleOut,
	}

	go rotLogger.monitorLogSize(staticFilePath)
	go rotLogger.processLogMessages()

	return rotLogger
}

func (rl *RotatingLogger) processLogMessages() {
	for logMessage := range rl.logChannel {
		if logMessage.Level >= rl.logLevel {
			entry := rl.Logger.WithFields(logMessage.Fields)
			switch logMessage.Level {
			case logrus.DebugLevel:
				entry.Debug(logMessage.Message)
			case logrus.InfoLevel:
				entry.Info(logMessage.Message)
			case logrus.WarnLevel:
				entry.Warn(logMessage.Message)
			case logrus.ErrorLevel:
				entry.Error(logMessage.Message)
			case logrus.FatalLevel:
				entry.Fatal(logMessage.Message)
			case logrus.PanicLevel:
				entry.Panic(logMessage.Message)
			}
			rl.multiWriter.Flush()
		}
	}
}

func (rl *RotatingLogger) Log(level logrus.Level, message string, fields logrus.Fields) {
	rl.logChannel <- LogMessage{
		Level:   level,
		Message: message,
		Fields:  fields,
	}
}

func (rl *RotatingLogger) monitorLogSize(staticFilePath string) {
	for {
		time.Sleep(rl.checkInterval)

		fileInfo, err := os.Stat(staticFilePath)
		if err != nil {
			rl.Logger.Errorf("Failed to get log file info: %v", err)
			continue
		}

		// Check if the file size exceeds the maxSize limit
		if fileInfo.Size() > int64(rl.maxSize*1024*1024) {
			now := time.Now()
			archiveFilename := filepath.Join(rl.logDir, fmt.Sprintf(rl.archivePattern, now.Format("2006-01-02-15-04-05")))
			rl.multiWriter.Flush()

			logrus.Infof("Log file size exceeded the limit. Renamed log file to: %s", archiveFilename)
			rl.Logger.Infof("Log file size exceeded the limit. Renamed log file to: "+archiveFilename, logrus.Fields{"appName": "MyApp"})

			err := os.Rename(staticFilePath, archiveFilename)
			if err != nil {
				rl.Logger.Errorf("Failed to rename log file: %v", err)
				continue
			}

			// Create a new log file after renaming the old one
			file, err := os.OpenFile(staticFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				rl.Logger.Errorf("Failed to create new log file: %v", err)
				continue
			}
			//rl.Logger.SetOutput(file)

			//multiWriter := io.Writer(file)
			multiWriter := bufio.NewWriter(file)
			//rl.Logger.SetOutput(multiWriter)

			if rl.consoleOut {
				fWriter := io.Writer(file)
				mWriter := io.MultiWriter(fWriter, os.Stdout)
				rl.Logger.SetOutput(mWriter)
			} else {
				rl.Logger.SetOutput(multiWriter)
			}

			rl.Logger.WithFields(logrus.Fields{"appName": "MyApp"}).Infof("Archive log file name: " + archiveFilename)

			if rl.zippedArchive {
				if err := zipFile(archiveFilename); err != nil {
					rl.Logger.Errorf("Failed to zip log file: %v", err)
					continue
				}
				if err := os.Remove(archiveFilename); err != nil {
					rl.Logger.Errorf("Failed to remove old log file: %v", err)
					continue
				}
			}

			rl.cleanupOldLogs()
		}
	}
}

func zipFile(source string) error {
	zipfile, err := os.Create(source + ".zip")
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	fileToZip, err := os.Open(source)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(source)
	header.Method = zip.Deflate

	writer, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}

func (rl *RotatingLogger) cleanupOldLogs() {
	files, err := filepath.Glob(filepath.Join(rl.logDir, "application-*.log*"))
	if err != nil {
		rl.Logger.Errorf("Failed to list log files: %v", err)
		return
	}

	sort.Slice(files, func(i, j int) bool {
		fi, _ := os.Stat(files[i])
		fj, _ := os.Stat(files[j])
		return fi.ModTime().Before(fj.ModTime())
	})

	if len(files) > rl.maxBackups {
		for _, file := range files[:len(files)-rl.maxBackups] {
			if err := os.Remove(file); err != nil {
				rl.Logger.Errorf("Failed to remove old log file: %v", err)
			} else {
				rl.Logger.WithFields(logrus.Fields{"appName": "MyApp"}).Infof("Removed old log file: %s", file)
			}
		}
	}
}

// JournalctlFormatter formats logs in a style similar to journalctl
type JournalctlFormatter struct{}

func (f *JournalctlFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("Jan 02 15:04:05")
	host, _ := os.Hostname()
	message := fmt.Sprintf("%s %s %s[%d]: %s\n", timestamp, host, entry.Data["appName"], os.Getpid(), entry.Message)
	return []byte(message), nil
}
