package configs_test

import (
	"testing"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig(t *testing.T) {
	// Helpers for options (reuse those from configs package for consistency)
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

	withStoreInterval := func(interval int) configs.ServerOption {
		return func(cfg *configs.ServerConfig) {
			cfg.StoreInterval = interval
		}
	}

	withFileStoragePath := func(path string) configs.ServerOption {
		return func(cfg *configs.ServerConfig) {
			cfg.FileStoragePath = path
		}
	}

	withRestore := func(restore bool) configs.ServerOption {
		return func(cfg *configs.ServerConfig) {
			cfg.Restore = restore
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
			want: &configs.ServerConfig{
				StoreInterval:   300,
				FileStoragePath: "metrics.json",
				Restore:         false,
			},
		},
		{
			name:    "set address only",
			options: []configs.ServerOption{withAddress("127.0.0.1:8080")},
			want: &configs.ServerConfig{
				Address:         "127.0.0.1:8080",
				StoreInterval:   300,
				FileStoragePath: "metrics.json",
				Restore:         false,
			},
		},
		{
			name:    "set log level only",
			options: []configs.ServerOption{withLogLevel("debug")},
			want: &configs.ServerConfig{
				LogLevel:        "debug",
				StoreInterval:   300,
				FileStoragePath: "metrics.json",
				Restore:         false,
			},
		},
		{
			name: "set all fields",
			options: []configs.ServerOption{
				withAddress("0.0.0.0:9000"),
				withLogLevel("info"),
				withStoreInterval(60),
				withFileStoragePath("/tmp/m.json"),
				withRestore(true),
			},
			want: &configs.ServerConfig{
				Address:         "0.0.0.0:9000",
				LogLevel:        "info",
				StoreInterval:   60,
				FileStoragePath: "/tmp/m.json",
				Restore:         true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := configs.NewServerConfig(tt.options...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWithServerAddress(t *testing.T) {
	cfg := &configs.ServerConfig{}
	opt := configs.WithServerAddress("127.0.0.1:8080")
	opt(cfg)
	assert.Equal(t, "127.0.0.1:8080", cfg.Address)
}

func TestWithServerLogLevel(t *testing.T) {
	cfg := &configs.ServerConfig{}
	opt := configs.WithServerLogLevel("debug")
	opt(cfg)
	assert.Equal(t, "debug", cfg.LogLevel)
}

func TestWithStoreInterval(t *testing.T) {
	cfg := &configs.ServerConfig{}
	opt := configs.WithStoreInterval(120)
	opt(cfg)
	assert.Equal(t, 120, cfg.StoreInterval)
}

func TestWithFileStoragePath(t *testing.T) {
	cfg := &configs.ServerConfig{}
	opt := configs.WithFileStoragePath("/tmp/metrics.json")
	opt(cfg)
	assert.Equal(t, "/tmp/metrics.json", cfg.FileStoragePath)
}

func TestWithRestore(t *testing.T) {
	cfg := &configs.ServerConfig{}
	opt := configs.WithRestore(true)
	opt(cfg)
	assert.True(t, cfg.Restore)

	optFalse := configs.WithRestore(false)
	optFalse(cfg)
	assert.False(t, cfg.Restore)
}
