package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	addressFlag        string
	pollIntervalFlag   int
	reportIntervalFlag int
	numWorkersFlag     int
)

func init() {
	parseAddressFlag()
	parsePollIntervalFlag()
	parseReportIntervalFlag()
	parseNumWorkersFlag()
}

func parseAddressFlag() {
	flag.StringVar(&addressFlag, "a", ":8080", "HTTP server address")
	if env := os.Getenv("ADDRESS"); env != "" {
		addressFlag = env
	}
}

func parsePollIntervalFlag() {
	flag.IntVar(&pollIntervalFlag, "poll", 5, "Polling interval in seconds")
	if env := os.Getenv("POLL_INTERVAL"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			pollIntervalFlag = v
		}
	}
}

func parseReportIntervalFlag() {
	flag.IntVar(&reportIntervalFlag, "report", 10, "Reporting interval in seconds")
	if env := os.Getenv("REPORT_INTERVAL"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			reportIntervalFlag = v
		}
	}
}

func parseNumWorkersFlag() {
	flag.IntVar(&numWorkersFlag, "workers", 4, "Number of worker goroutines")
	if env := os.Getenv("NUM_WORKERS"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			numWorkersFlag = v
		}
	}
}
