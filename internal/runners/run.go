// Package runners provides utilities to manage the lifecycle of long-running services or agents.
package runners

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// Runnable defines an interface for components that have a start/stop lifecycle,
// such as servers, workers, or agents.
type Runnable interface {
	// Start launches the runnable component. It should block until the component
	// exits or an error occurs.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the component. It receives a context that can be
	// used to control the shutdown timeout or cancellation.
	Stop(ctx context.Context) error
}

// Run executes a Runnable component and manages its lifecycle.
//
// It listens for termination signals (SIGINT, SIGTERM, SIGQUIT) and cancels
// the context when such a signal is received. Upon cancellation, Run attempts
// a graceful shutdown by calling Stop on the Runnable with a 5-second timeout.
//
// Run also captures errors returned by Start. If Start returns an error other
// than http.ErrServerClosed, Run returns that error immediately.
//
// Run blocks until either the Runnable completes or a termination signal is received.
//
// Parameters:
//   - ctx: The parent context to control the lifecycle.
//   - runnable: The Runnable instance to start and stop.
//
// Returns:
//   - An error if the Runnable fails to start or stops with an error (other than http.ErrServerClosed).
//
// Usage example:
//
//	err := runners.Run(ctx, myApp)
//	if err != nil {
//	    log.Fatal(err)
//	}
func Run(
	ctx context.Context,
	runnable Runnable,
) error {
	ctx, cancel := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		err := runnable.Start(ctx)
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return runnable.Stop(shutdownCtx)

	case err := <-errChan:
		return err
	}
}
