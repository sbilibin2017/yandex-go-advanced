package main

import (
	"context"
	"testing"
	"time"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	// Create a context with timeout so the server stops automatically
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create a minimal config, adjust fields as needed
	cfg := &configs.ServerConfig{
		Address:  "localhost:0", // 0 means random free port, if supported
		LogLevel: "info",
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- run(ctx, cfg)
	}()

	// Wait for run to exit after context timeout
	err := <-errCh

	// On graceful shutdown, expect no error (or context canceled depending on implementation)
	assert.NoError(t, err)
}
