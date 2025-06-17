package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestParseFlags_EnvOverridesFlags(t *testing.T) {
	tests := []struct {
		name         string
		env          map[string]string
		args         []string
		wantAddr     string
		wantPoll     int
		wantReport   int
		wantWorkers  int
		wantLogLevel string
	}{
		{
			name: "env overrides flags",
			env: map[string]string{
				"ADDRESS":         "envhost:9090",
				"POLL_INTERVAL":   "5",
				"REPORT_INTERVAL": "20",
				"NUM_WORKERS":     "8",
				"LOG_LEVEL":       "debug",
			},
			args: []string{"cmd",
				"-a", "flaghost:7070",
				"-p", "15",
				"-r", "25",
				"-workers", "16",
				"-l", "warn",
			},
			wantAddr:     "envhost:9090",
			wantPoll:     5,
			wantReport:   20,
			wantWorkers:  8,
			wantLogLevel: "debug",
		},
		{
			name: "flags only",
			env:  nil,
			args: []string{"cmd",
				"-a", "flaghost:7070",
				"-p", "15",
				"-r", "25",
				"-workers", "16",
				"-l", "warn",
			},
			wantAddr:     "flaghost:7070",
			wantPoll:     15,
			wantReport:   25,
			wantWorkers:  16,
			wantLogLevel: "warn",
		},
		{
			name: "env only",
			env: map[string]string{
				"ADDRESS":         "envhost:9090",
				"POLL_INTERVAL":   "5",
				"REPORT_INTERVAL": "20",
				"NUM_WORKERS":     "8",
				"LOG_LEVEL":       "debug",
			},
			args:         []string{"cmd"},
			wantAddr:     "envhost:9090",
			wantPoll:     5,
			wantReport:   20,
			wantWorkers:  8,
			wantLogLevel: "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env first
			for k := range tt.env {
				os.Unsetenv(k)
			}
			// Set env vars for this test
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			resetFlags()
			os.Args = tt.args

			serverAddr = ""
			pollInterval = 0
			reportInterval = 0
			numWorkers = 0
			logLevel = ""

			parseFlags()

			assert.Equal(t, tt.wantAddr, serverAddr)
			assert.Equal(t, tt.wantPoll, pollInterval)
			assert.Equal(t, tt.wantReport, reportInterval)
			assert.Equal(t, tt.wantWorkers, numWorkers)
			assert.Equal(t, tt.wantLogLevel, logLevel)

			for k := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}
