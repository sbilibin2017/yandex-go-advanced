package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	// Set up flags or variables used by run
	addr = "localhost:0" // use port 0 for automatic free port
	logLevel = "debug"

	// Create a context with timeout to stop run gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := run(ctx)
	assert.NoError(t, err)
}
