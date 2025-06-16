package configs_test

import (
	"testing"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig(t *testing.T) {
	// Option helpers
	withAddress := func(addr string) configs.ServerOption {
		return func(cfg *configs.ServerConfig) {
			cfg.Address = addr
		}
	}

	withLogLevel := func(level string) configs.ServerOption {
		return func(cfg *configs.ServerConfig) {
			cfg.LogLevel = level
		}
	}

	tests := []struct {
		name    string
		options []configs.ServerOption
		want    *configs.ServerConfig
	}{
		{
			name:    "default config",
			options: nil,
			want:    &configs.ServerConfig{},
		},
		{
			name:    "set address only",
			options: []configs.ServerOption{withAddress("127.0.0.1:8080")},
			want:    &configs.ServerConfig{Address: "127.0.0.1:8080"},
		},
		{
			name:    "set log level only",
			options: []configs.ServerOption{withLogLevel("debug")},
			want:    &configs.ServerConfig{LogLevel: "debug"},
		},
		{
			name:    "set address and log level",
			options: []configs.ServerOption{withAddress("0.0.0.0:9000"), withLogLevel("info")},
			want:    &configs.ServerConfig{Address: "0.0.0.0:9000", LogLevel: "info"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := configs.NewServerConfig(tt.options...)
			assert.Equal(t, tt.want, got)
		})
	}
}
