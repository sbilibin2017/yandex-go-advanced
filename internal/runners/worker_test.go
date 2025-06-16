package runners

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestRunWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	started := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)

	worker := func(workerCtx context.Context) {
		defer wg.Done()

		// Signal worker started
		started <- struct{}{}

		// Wait for context to be canceled
		<-workerCtx.Done()
	}

	go RunWorker(ctx, worker)

	// Wait until worker starts
	<-started

	// Cancel the context, which should cause RunWorker to unblock and return
	cancel()

	// Give a moment for the worker to finish after context cancellation
	time.Sleep(50 * time.Millisecond)

	wg.Wait()
}
