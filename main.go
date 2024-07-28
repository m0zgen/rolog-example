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
			s := strconv.Itoa(i)
			ro.Logger.Log(logrus.InfoLevel, "This is an info log message. With item: "+s, logrus.Fields{"appName": "exampleApp"})
		}

		time.Sleep(1 * time.Second)
	}

}

// GenerateChannelLogs - Generate logs using channel
func GenerateChannelLogs(id int, data chan int) {
	for taskId := range data {
		sId := strconv.Itoa(id)
		ro.Logger.Log(logrus.InfoLevel, "This is an info log message from channel. With worker and id: "+sId+"-"+strconv.Itoa(taskId), logrus.Fields{"appName": "exampleApp"})

		time.Sleep(2 * time.Second)
	}
}

// RunTimeTest - Run time test
func RunTimeTest() {
	var totalDuration time.Duration
	runs := 100
	messagesPerRun := 1000

	for i := 0; i < runs; i++ {
		start := time.Now()

		for j := 0; j < messagesPerRun; j++ {
			ro.Logger.Log(logrus.InfoLevel, "This is an info log message.", logrus.Fields{"appName": "MyApp"})
		}

		duration := time.Since(start)
		totalDuration += duration
		logrus.Infof("Run %d: Time taken to log %d messages: %v", i+1, messagesPerRun, duration)
	}

	averageDuration := totalDuration / time.Duration(runs)
	logrus.Infof("Average time taken to log %d messages over %d runs: %v", messagesPerRun, runs, averageDuration)

}

func main() {

	RunTimeTest()

	// Create channel
	channel := make(chan int)

	// Creating 100 workers to execute task
	for i := 0; i < 1; i++ {
		go GenerateChannelLogs(i, channel)
	}

	// Filling channel with 100000 data
	for i := 0; i < 10; i++ {
		channel <- i
	}

	// Simple log generation
	GenerateLogs()

}
