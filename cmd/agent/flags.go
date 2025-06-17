package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	serverAddr     string
	pollInterval   int
	reportInterval int
	numWorkers     int
	logLevel       string
)

func parseFlags() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "server address")
	flag.IntVar(&pollInterval, "p", 2, "polling interval in seconds")
	flag.IntVar(&reportInterval, "r", 10, "reporting interval in seconds")
	flag.IntVar(&numWorkers, "workers", 4, "number of workers")
	flag.StringVar(&logLevel, "l", "info", "log level")

	flag.Parse()

	if env := os.Getenv("ADDRESS"); env != "" {
		serverAddr = env
	}
	if env := os.Getenv("POLL_INTERVAL"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			pollInterval = v
		}
	}
	if env := os.Getenv("REPORT_INTERVAL"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			reportInterval = v
		}
	}
	if env := os.Getenv("NUM_WORKERS"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			numWorkers = v
		}
	}
	if env := os.Getenv("LOG_LEVEL"); env != "" {
		logLevel = env
	}
}
