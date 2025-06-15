package main

import (
	"flag"
	"os"
)

var (
	addressFlag string
)

func init() {
	parseAddressFlag()
}

func parseAddressFlag() {
	flag.StringVar(&addressFlag, "a", ":8080", "HTTP server address")
	if env := os.Getenv("ADDRESS"); env != "" {
		addressFlag = env
	}
}
