package main

import (
	"rolog-example/internal/rotatinglogger"
	"time"
)

func main() {
	// Настройки логирования
	filename := "application-%s.log"
	datePattern := "2006-01-02-15"
	zippedArchive := true
	maxSize := 20  // 20 MB
	maxFiles := 14 // 14 days

	logger := rotatinglogger.NewRotatingLogger(filename, datePattern, zippedArchive, maxSize, maxFiles)

	for {
		logger.Logger.Info("This is a log message.")
		time.Sleep(1 * time.Hour)
	}
}
