package apps

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestServerApp_StartAndStop(t *testing.T) {
	// Create config with a free port for testing
	cfg := &configs.ServerConfig{
		Address: "127.0.0.1:0", // port 0 means assign any available port
	}

	app, err := NewServerApp(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	// Start server in background goroutine because ListenAndServe blocks
	go func() {
		err := app.Start(context.Background())
		assert.ErrorIs(t, err, http.ErrServerClosed, "expected ErrServerClosed on shutdown")
	}()

	// Wait briefly to let server start
	time.Sleep(100 * time.Millisecond)

	// Create a context with timeout to stop the server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Stop the server and assert no error
	err = app.Stop(ctx)
	assert.NoError(t, err)
}

func TestServerWorkerApp_StartAndStop(t *testing.T) {
	// Prepare a test config
	cfg := &configs.ServerConfig{
		FileStoragePath: "testdata/metrics.json",
		StoreInterval:   1,     // 1 second interval for quick test
		Restore:         false, // do not restore to simplify test
	}

	workerApp, err := NewServerWorkerApp(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, workerApp)

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Run Start in a separate goroutine since it blocks
	done := make(chan struct{})
	go func() {
		err := workerApp.Start(ctx)
		assert.NoError(t, err) // Start always returns nil
		close(done)
	}()

	// Wait a short moment to simulate some work done by the worker
	time.Sleep(100 * time.Millisecond)

	// Cancel the context to stop the worker
	cancel()

	// Wait for the Start goroutine to finish
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for worker to stop")
	}

	// Stopping should not produce errors
	err = workerApp.Stop(context.Background())
	assert.NoError(t, err)
}
