package ro

import (
	"github.com/sirupsen/logrus"
	"rolog-example/internal/rotatinglogger"
	"time"
)

var Logger *rotatinglogger.RotatingLogger

func init() {
	// Logger Settings
	logDir := "logs"
	staticFilename := "application.log"
	archivePattern := "application-%s.log"
	zippedArchive := true        // Enable zipped archive
	maxSize := 20                // 20 MB
	maxBackups := 3              // 3 backups
	checkInterval := time.Hour   // Check every hour
	bufferSize := 100            // Size of the log message buffer
	logLevel := logrus.InfoLevel // Set the log level to Info

	Logger = rotatinglogger.NewRotatingLogger(
		logDir,
		staticFilename,
		archivePattern,
		zippedArchive,
		maxSize,
		maxBackups,
		checkInterval,
		bufferSize,
		logLevel,
	)
}
