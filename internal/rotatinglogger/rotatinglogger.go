package rotatinglogger

import (
	"archive/zip"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RotatingLogger struct {
	Logger *logrus.Logger
}

func NewRotatingLogger(filename string, datePattern string, zippedArchive bool, maxSize int, maxFiles int) *RotatingLogger {
	logger := logrus.New()

	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	// Generate the log file name based on the date pattern
	logFileName := generateLogFileName(filename, datePattern)

	logger.Out = &lumberjack.Logger{
		Filename: logFileName,
		MaxSize:  maxSize,  // megabytes
		MaxAge:   maxFiles, // days
		Compress: zippedArchive,
	}

	return &RotatingLogger{
		Logger: logger,
	}
}

func generateLogFileName(filename string, datePattern string) string {
	now := time.Now()
	dateFormatted := now.Format(convertDatePattern(datePattern))
	return fmt.Sprintf(filename, dateFormatted)
}

func convertDatePattern(datePattern string) string {
	replacer := strings.NewReplacer(
		"YYYY", "2006",
		"MM", "01",
		"DD", "02",
		"HH", "15",
		"mm", "04",
		"ss", "05",
	)
	return replacer.Replace(datePattern)
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
