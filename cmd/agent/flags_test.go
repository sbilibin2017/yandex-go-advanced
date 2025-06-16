package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags_TableDriven(t *testing.T) {
	tests := []struct {
		name            string
		envVars         map[string]string
		args            []string
		expectedAddress string
		expectedPoll    int
		expectedReport  int
		expectedWorkers int
		expectedLog     string
	}{
		{
			name:            "defaults",
			envVars:         nil,
			args:            []string{"cmd"},
			expectedAddress: "localhost:8080",
			expectedPoll:    2,
			expectedReport:  10,
			expectedWorkers: 4,
			expectedLog:     "info",
		},
		{
			name:    "flags override defaults",
			envVars: nil,
			args: []string{
				"cmd",
				"-a=flagaddress:1234",
				"-p=7",
				"-r=20",
				"-workers=9",
				"-l=debug",
			},
			expectedAddress: "flagaddress:1234",
			expectedPoll:    7,
			expectedReport:  20,
			expectedWorkers: 9,
			expectedLog:     "debug",
		},
		{
			name: "env overrides flags",
			envVars: map[string]string{
				"ADDRESS":         "envaddress:4321",
				"POLL_INTERVAL":   "15",
				"REPORT_INTERVAL": "30",
				"NUM_WORKERS":     "8",
				"LOG_LEVEL":       "warn",
			},
			args: []string{
				"cmd",
				"-a=flagaddress",
				"-p=7",
				"-r=20",
				"-workers=9",
				"-l=debug",
			},
			expectedAddress: "envaddress:4321",
			expectedPoll:    15,
			expectedReport:  30,
			expectedWorkers: 8,
			expectedLog:     "warn",
		},
		{
			name: "invalid env falls back to flags",
			envVars: map[string]string{
				"POLL_INTERVAL":   "notanint",
				"REPORT_INTERVAL": "alsobad",
				"NUM_WORKERS":     "badnum",
			},
			args: []string{
				"cmd",
				"-p=5",
				"-r=25",
				"-workers=3",
			},
			expectedAddress: "localhost:8080",
			expectedPoll:    5,
			expectedReport:  25,
			expectedWorkers: 3,
			expectedLog:     "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnvVars(t)
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			os.Args = tt.args

			cfg, err := parseFlags()
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedAddress, cfg.ServerAddress)
			assert.Equal(t, tt.expectedPoll, cfg.PollInterval)
			assert.Equal(t, tt.expectedReport, cfg.ReportInterval)
			assert.Equal(t, tt.expectedWorkers, cfg.NumWorkers)
			assert.Equal(t, tt.expectedLog, cfg.LogLevel)

			clearEnvVars(t)
		})
	}
}

func clearEnvVars(t *testing.T) {
	t.Helper()
	os.Unsetenv("ADDRESS")
	os.Unsetenv("POLL_INTERVAL")
	os.Unsetenv("REPORT_INTERVAL")
	os.Unsetenv("NUM_WORKERS")
	os.Unsetenv("LOG_LEVEL")
}
