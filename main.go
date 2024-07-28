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

func main() {

	GenerateLogs()

}
