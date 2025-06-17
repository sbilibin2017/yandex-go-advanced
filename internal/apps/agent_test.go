package apps

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestNewAgentApp(t *testing.T) {
	cfg := &configs.AgentConfig{
		ServerAddress:  "http://localhost:8080",
		PollInterval:   1000,
		ReportInterval: 2000,
	}

	app, err := NewAgentApp(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.NotNil(t, app.worker, "worker func should not be nil")
}

func TestAgentApp_StartAndStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup a mock worker function which respects context cancellation
	var wg sync.WaitGroup
	wg.Add(1)

	mockWorker := func(ctx context.Context) {
		defer wg.Done()
		<-ctx.Done() // block until context canceled
	}

	app := &AgentApp{
		worker: mockWorker,
	}

	// Start the worker in a separate goroutine so test can continue
	go func() {
		err := app.Start(ctx)
		assert.NoError(t, err)
	}()

	// Give the worker some time to start
	time.Sleep(50 * time.Millisecond)

	// Stop by canceling context, which should cause worker to exit
	cancel()

	// Wait for worker to finish
	wg.Wait()

	// Stop currently does nothing but test it anyway
	err := app.Stop(context.Background())
	assert.NoError(t, err)
}
