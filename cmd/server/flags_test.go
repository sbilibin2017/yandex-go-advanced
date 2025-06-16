package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags_ValidInputs(t *testing.T) {
	tests := []struct {
		name            string
		envVars         map[string]string
		args            []string
		expectedAddress string
		expectedLog     string
	}{
		{
			name:            "default values",
			envVars:         nil,
			args:            []string{"cmd"},
			expectedAddress: "localhost:8080",
			expectedLog:     "info",
		},
		{
			name:    "flags override defaults",
			envVars: nil,
			args: []string{
				"cmd",
				"-a=flagaddress:1234",
				"-l=debug",
			},
			expectedAddress: "flagaddress:1234",
			expectedLog:     "debug",
		},
		{
			name: "env overrides flags",
			envVars: map[string]string{
				"ADDRESS":   "envaddress:4321",
				"LOG_LEVEL": "warn",
			},
			args: []string{
				"cmd",
				"-a=flagaddress",
				"-l=debug",
			},
			expectedAddress: "envaddress:4321",
			expectedLog:     "warn",
		},
		{
			name: "env unset falls back to flags",
			envVars: map[string]string{
				"ADDRESS":   "",
				"LOG_LEVEL": "",
			},
			args: []string{
				"cmd",
				"-a=flagaddress",
				"-l=debug",
			},
			expectedAddress: "flagaddress",
			expectedLog:     "debug",
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

			assert.Equal(t, tt.expectedAddress, cfg.Address)
			assert.Equal(t, tt.expectedLog, cfg.LogLevel)

			clearEnvVars(t)
		})
	}
}

func clearEnvVars(t *testing.T) {
	t.Helper()
	os.Unsetenv("ADDRESS")
	os.Unsetenv("LOG_LEVEL")
}
