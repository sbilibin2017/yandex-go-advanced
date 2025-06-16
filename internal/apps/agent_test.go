package apps

import (
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

	assert.NotNil(t, worker)
	assert.Nil(t, err)

}
