package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestRun_CustomConfig_WithTestServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Addr:    "localhost:1234",
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	time.Sleep(100 * time.Millisecond)

	cfg := configs.NewAgentConfig(func(c *configs.AgentConfig) {
		c.ServerAddress = "localhost:1234"
		c.PollInterval = 1
		c.ReportInterval = 5
		c.NumWorkers = 1
		c.LogLevel = "debug"
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	doneCh := make(chan error)
	go func() {
		doneCh <- run(ctx, cfg)
	}()

	time.Sleep(500 * time.Millisecond)

	cancel() // simulate early shutdown

	err := <-doneCh

	// Accept context cancellation as non-error
	if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		t.Fatalf("unexpected error: %v", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdownCancel()
	err = srv.Shutdown(shutdownCtx)
	assert.NoError(t, err)

	if serveErr := <-errCh; serveErr != nil {
		t.Fatalf("server error: %v", serveErr)
	}
}
