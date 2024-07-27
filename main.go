package main

import (
	"github.com/sirupsen/logrus"
	"rolog-example/internal/rotatinglogger"
	"time"
)

func main() {
	// Настройки логирования
	logDir := "logs"
	staticFilename := "application.log"
	archivePattern := "application-%s.log"
	zippedArchive := true
	maxSize := 20              // 20 MB
	maxAge := 14               // 14 days
	maxBackups := 3            // 3 backups
	checkInterval := time.Hour // Check every hour
	bufferSize := 100          // Size of the log message buffer

	logger := rotatinglogger.NewRotatingLogger(logDir, staticFilename, archivePattern, zippedArchive, maxSize, maxAge, maxBackups, checkInterval, bufferSize)

	for {
		logger.Log(logrus.InfoLevel, "This is a log message.", logrus.Fields{"appName": "exampleApp"})
		time.Sleep(1 * time.Hour)
	}
}
