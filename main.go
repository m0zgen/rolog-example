package main

import (
	"rolog-example/internal/rotatinglogger"
	"time"
)

func main() {
	// Настройки логирования
	logDir := "logs"
	staticFilename := "application.log"
	archivePattern := "application-%s.log"
	zippedArchive := true
	maxSize := 1               // 20 MB
	maxAge := 14               // 14 days
	maxBackups := 3            // 3 backups
	checkInterval := time.Hour // Check every hour

	logger := rotatinglogger.NewRotatingLogger(logDir, staticFilename, archivePattern, zippedArchive, maxSize, maxAge, maxBackups, checkInterval)

	for {
		logger.Logger.WithField("appName", "exampleApp").Info("This is a log message.")
		time.Sleep(1 * time.Second)
	}
}
