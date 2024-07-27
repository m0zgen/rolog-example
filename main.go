package main

import (
	"github.com/sirupsen/logrus"
	"rolog-example/internal/rotatinglogger"
	"strconv"
	"time"
)

func main() {
	// Logger Settings
	logDir := "logs"
	staticFilename := "application.log"
	archivePattern := "application-%s.log"
	zippedArchive := true        // Enable zipped archive
	maxSize := 1                 // 20 MB
	maxBackups := 3              // 3 backups
	checkInterval := time.Second // Check every hour
	bufferSize := 100            // Size of the log message buffer
	logLevel := logrus.InfoLevel // Set the log level to Info

	//logger := rotatinglogger.NewRotatingLogger(logDir, staticFilename, archivePattern, zippedArchive, maxSize, maxBackups, checkInterval, bufferSize, logLevel)

	logger := rotatinglogger.NewRotatingLogger(
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

	for {
		logger.Log(logrus.InfoLevel, "This is an info log message.", logrus.Fields{"appName": "exampleApp"})
		logger.Log(logrus.DebugLevel, "This is a debug log message.", logrus.Fields{"appName": "exampleApp"})

		// Generate some log messages - 1000 messages
		for i := 0; i < 1000; i++ {
			// i to string
			s := strconv.Itoa(i)
			logger.Log(logrus.InfoLevel, "This is an info log message. Number: "+s, logrus.Fields{"appName": "exampleApp"})
		}

		time.Sleep(1 * time.Second)
	}
}
