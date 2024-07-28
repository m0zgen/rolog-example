package main

import (
	"github.com/sirupsen/logrus"
	"rolog-example/internal/ro"
	"strconv"
	"time"
)

func GenerateLogs() {

	for {

		// Generate some log messages - 1000 messages
		for i := 0; i < 1000; i++ {
			// i to string
			s := strconv.Itoa(i)
			ro.Logger.Log(logrus.InfoLevel, "This is an info log message. With item: "+s, logrus.Fields{"appName": "exampleApp"})
		}

		time.Sleep(1 * time.Second)
	}

}

// GenerateChannelLogs - Generate logs using channel
func GenerateChannelLogs(id int, data chan int) {
	for taskId := range data {
		time.Sleep(2 * time.Second)
		ro.Logger.Log(logrus.InfoLevel, "This is an info log message from channel. With item: "+strconv.Itoa(taskId), logrus.Fields{"appName": "exampleApp"})
	}
}

func main() {

	// Create channel
	channel := make(chan int)

	// Creating 100 workers to execute task
	for i := 0; i < 100; i++ {
		go GenerateChannelLogs(i, channel)
	}

	// Filling channel with data
	for i := 0; i < 100000; i++ {
		channel <- i
	}

	// Simple log generation
	GenerateLogs()

}
