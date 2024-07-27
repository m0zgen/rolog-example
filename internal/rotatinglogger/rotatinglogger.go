package rotatinglogger

import (
	"archive/zip"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
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
}

func NewRotatingLogger(logDir, staticFilename, archivePattern string, zippedArchive bool, maxSize, maxAge, maxBackups int, checkInterval time.Duration, bufferSize int) *RotatingLogger {
	logger := logrus.New()

	logger.Formatter = &JournalctlFormatter{}

	// Ensure log directory exists
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		logger.Fatalf("Failed to create log directory: %v", err)
	}

	staticFilePath := filepath.Join(logDir, staticFilename)

	logger.Out = &lumberjack.Logger{
		Filename:   staticFilePath,
		MaxSize:    maxSize,    // megabytes
		MaxAge:     maxAge,     // days
		MaxBackups: maxBackups, // backups
		Compress:   false,      // Disable built-in compression, we'll handle it
	}

	rotLogger := &RotatingLogger{
		Logger:         logger,
		checkInterval:  checkInterval,
		maxSize:        maxSize,
		zippedArchive:  zippedArchive,
		logDir:         logDir,
		archivePattern: archivePattern,
		logChannel:     make(chan LogMessage, bufferSize),
	}

	go rotLogger.monitorLogSize(staticFilePath)
	go rotLogger.processLogMessages()

	return rotLogger
}

func (rl *RotatingLogger) processLogMessages() {
	for logMessage := range rl.logChannel {
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
			err := os.Rename(staticFilePath, archiveFilename)
			if err != nil {
				rl.Logger.Errorf("Failed to rename log file: %v", err)
				continue
			}
			if rl.zippedArchive {
				if err := zipFile(archiveFilename); err != nil {
					rl.Logger.Errorf("Failed to zip log file: %v", err)
					continue
				}
				if err := os.Remove(archiveFilename); err != nil {
					rl.Logger.Errorf("Failed to remove old log file: %v", err)
				}
			}

			// Create a new log file after renaming the old one
			_, err = os.Create(staticFilePath)
			if err != nil {
				rl.Logger.Errorf("Failed to create new log file: %v", err)
			}
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

// JournalctlFormatter formats logs in a style similar to journalctl
type JournalctlFormatter struct{}

func (f *JournalctlFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("Jan 02 15:04:05")
	host, _ := os.Hostname()
	message := fmt.Sprintf("%s %s %s[%d]: %s\n", timestamp, host, entry.Data["appName"], os.Getpid(), entry.Message)
	return []byte(message), nil
}
