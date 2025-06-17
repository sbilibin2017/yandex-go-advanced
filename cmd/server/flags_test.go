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

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name         string
		env          map[string]string
		args         []string
		wantAddr     string
		wantLogLevel string
		wantInterval int
		wantFilePath string
		wantRestore  bool
	}{
		{
			name: "env overrides flags",
			env: map[string]string{
				"ADDRESS":           "envhost:9090",
				"LOG_LEVEL":         "debug",
				"STORE_INTERVAL":    "123",
				"FILE_STORAGE_PATH": "/tmp/env_metrics.json",
				"RESTORE":           "true",
			},
			args:         []string{"cmd", "-a", "flaghost:7070", "-l", "warn", "-i", "456", "-f", "/tmp/flag_metrics.json", "-r", "false"},
			wantAddr:     "envhost:9090",
			wantLogLevel: "debug",
			wantInterval: 123,
			wantFilePath: "/tmp/env_metrics.json",
			wantRestore:  true,
		},
		{
			name:         "flags only",
			env:          nil,
			args:         []string{"cmd", "-a", "flaghost:7070", "-l", "warn", "-i", "456", "-f", "/tmp/flag_metrics.json", "-r", "true"},
			wantAddr:     "flaghost:7070",
			wantLogLevel: "warn",
			wantInterval: 456,
			wantFilePath: "/tmp/flag_metrics.json",
			wantRestore:  true,
		},
		{
			name: "env only",
			env: map[string]string{
				"ADDRESS":           "envhost:9090",
				"LOG_LEVEL":         "debug",
				"STORE_INTERVAL":    "123",
				"FILE_STORAGE_PATH": "/tmp/env_metrics.json",
				"RESTORE":           "true",
			},
			args:         []string{"cmd"},
			wantAddr:     "envhost:9090",
			wantLogLevel: "debug",
			wantInterval: 123,
			wantFilePath: "/tmp/env_metrics.json",
			wantRestore:  true,
		},
		{
			name:         "defaults without env or flags",
			env:          nil,
			args:         []string{"cmd"},
			wantAddr:     "localhost:8080",
			wantLogLevel: "info",
			wantInterval: 300,
			wantFilePath: "metrics_storage.json",
			wantRestore:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all relevant env vars first
			os.Unsetenv("ADDRESS")
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("STORE_INTERVAL")
			os.Unsetenv("FILE_STORAGE_PATH")
			os.Unsetenv("RESTORE")

			// Set env vars for test case
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			// Reset flags and args
			resetFlags()
			os.Args = tt.args

			// Reset global vars to zero values before parsing
			addr = ""
			logLevel = ""
			storeInterval = 0
			filePath = ""
			restore = false

			parseFlags()

			assert.Equal(t, tt.wantAddr, addr)
			assert.Equal(t, tt.wantLogLevel, logLevel)
			assert.Equal(t, tt.wantInterval, storeInterval)
			assert.Equal(t, tt.wantFilePath, filePath)
			assert.Equal(t, tt.wantRestore, restore)

			// Clean up env vars after test
			for k := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}
