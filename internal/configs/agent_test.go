package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentOption_ServerAddress(t *testing.T) {
	expected := "http://localhost:8080"
	opt := func(cfg *AgentConfig) {
		cfg.ServerAddress = expected
	}

	cfg := NewAgentConfig(opt)
	assert.Equal(t, expected, cfg.ServerAddress)
}

func TestAgentOption_ServerEndpoint(t *testing.T) {
	expected := "/metrics"
	opt := func(cfg *AgentConfig) {
		cfg.ServerEndpoint = expected
	}

	cfg := NewAgentConfig(opt)
	assert.Equal(t, expected, cfg.ServerEndpoint)
}

func TestAgentOption_LogLevel(t *testing.T) {
	expected := "debug"
	opt := func(cfg *AgentConfig) {
		cfg.LogLevel = expected
	}

	cfg := NewAgentConfig(opt)
	assert.Equal(t, expected, cfg.LogLevel)
}

func TestAgentOption_PollInterval(t *testing.T) {
	expected := 15
	opt := func(cfg *AgentConfig) {
		cfg.PollInterval = expected
	}

	cfg := NewAgentConfig(opt)
	assert.Equal(t, expected, cfg.PollInterval)
}

func TestAgentOption_ReportInterval(t *testing.T) {
	expected := 30
	opt := func(cfg *AgentConfig) {
		cfg.ReportInterval = expected
	}

	cfg := NewAgentConfig(opt)
	assert.Equal(t, expected, cfg.ReportInterval)
}

func TestAgentOption_NumWorkers(t *testing.T) {
	expected := 4
	opt := func(cfg *AgentConfig) {
		cfg.NumWorkers = expected
	}

	cfg := NewAgentConfig(opt)
	assert.Equal(t, expected, cfg.NumWorkers)
}
