package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	addr          string
	logLevel      string
	storeInterval int
	filePath      string
	restore       bool
)

func parseFlags() {
	flag.StringVar(&addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.IntVar(&storeInterval, "i", 300, "interval in seconds to store metrics to disk (0 = sync)")
	flag.StringVar(&filePath, "f", "metrics_storage.json", "file path to store metrics")
	flag.BoolVar(&restore, "r", false, "restore metrics from file on start (true/false)")

	flag.Parse()

	if env := os.Getenv("ADDRESS"); env != "" {
		addr = env
	}
	if env := os.Getenv("LOG_LEVEL"); env != "" {
		logLevel = env
	}
	if env := os.Getenv("STORE_INTERVAL"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			storeInterval = v
		}
	}
	if env := os.Getenv("FILE_STORAGE_PATH"); env != "" {
		filePath = env
	}
	if env := os.Getenv("RESTORE"); env != "" {
		if v, err := strconv.ParseBool(env); err == nil {
			restore = v
		}
	}
}
