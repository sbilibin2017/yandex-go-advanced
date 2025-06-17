package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	// Set up the global variables (or use your config vars)
	serverAddr = "localhost:0" // port 0 means OS assigns a free port
	pollInterval = 1           // small interval for fast test
	reportInterval = 1         // small interval for fast test
	numWorkers = 1
	logLevel = "debug"

	// Create a context with timeout to stop run gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := run(ctx)
	assert.NoError(t, err)
}
