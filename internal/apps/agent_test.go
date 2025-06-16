package apps

import (
	"context"
	"testing"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestNewAgentApp(t *testing.T) {
	config := &configs.AgentConfig{
		ServerAddress:  "http://localhost:8080",
		PollInterval:   1,
		ReportInterval: 1,
		NumWorkers:     2,
	}

	worker, err := NewAgentApp(config)

	assert.NoError(t, err)
	assert.NotNil(t, worker)

	// Optional: test that worker function runs without error for a short time with canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = worker(ctx)
	assert.ErrorIs(t, err, context.Canceled)
}
