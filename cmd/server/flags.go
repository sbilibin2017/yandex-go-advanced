package main

import (
	"flag"
	"os"
)

var (
	addr     string
	logLevel string
)

func parseFlags() {
	flag.StringVar(&addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&logLevel, "l", "info", "log level")

	flag.Parse()

	if env := os.Getenv("ADDRESS"); env != "" {
		addr = env
	}
	if env := os.Getenv("LOG_LEVEL"); env != "" {
		logLevel = env
	}
}
