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
		name       string
		env        map[string]string
		args       []string
		wantAddr   string
		wantLogLvl string
	}{
		{
			name: "env overrides flags",
			env: map[string]string{
				"ADDRESS":   "envhost:9090",
				"LOG_LEVEL": "debug",
			},
			args:       []string{"cmd", "-a", "flaghost:7070", "-l", "warn"},
			wantAddr:   "envhost:9090",
			wantLogLvl: "debug",
		},
		{
			name:       "flags only",
			env:        nil,
			args:       []string{"cmd", "-a", "flaghost:7070", "-l", "warn"},
			wantAddr:   "flaghost:7070",
			wantLogLvl: "warn",
		},
		{
			name: "env only",
			env: map[string]string{
				"ADDRESS":   "envhost:9090",
				"LOG_LEVEL": "debug",
			},
			args:       []string{"cmd"},
			wantAddr:   "envhost:9090",
			wantLogLvl: "debug",
		},
		{
			name:       "defaults without env or flags",
			env:        nil,
			args:       []string{"cmd"},
			wantAddr:   "localhost:8080",
			wantLogLvl: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env first
			os.Unsetenv("ADDRESS")
			os.Unsetenv("LOG_LEVEL")

			// Set env vars for test
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			resetFlags()
			os.Args = tt.args

			// Reset globals before parsing
			addr = ""
			logLevel = ""

			parseFlags()

			assert.Equal(t, tt.wantAddr, addr)
			assert.Equal(t, tt.wantLogLvl, logLevel)

			// Clean up env
			for k := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}
