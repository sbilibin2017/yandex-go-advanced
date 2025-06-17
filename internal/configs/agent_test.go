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

func TestWithAgentServerAddress(t *testing.T) {
	cfg := &AgentConfig{}
	opt := WithAgentServerAddress("localhost:8080")
	opt(cfg)
	assert.Equal(t, "localhost:8080", cfg.ServerAddress)
}

func TestWithAgentLogLevel(t *testing.T) {
	cfg := &AgentConfig{}
	opt := WithAgentLogLevel("info")
	opt(cfg)
	assert.Equal(t, "info", cfg.LogLevel)
}

func TestWithAgentPollInterval(t *testing.T) {
	cfg := &AgentConfig{}
	opt := WithAgentPollInterval(10)
	opt(cfg)
	assert.Equal(t, 10, cfg.PollInterval)
}

func TestWithAgentReportInterval(t *testing.T) {
	cfg := &AgentConfig{}
	opt := WithAgentReportInterval(20)
	opt(cfg)
	assert.Equal(t, 20, cfg.ReportInterval)
}

func TestWithAgentNumWorkers(t *testing.T) {
	cfg := &AgentConfig{}
	opt := WithAgentNumWorkers(5)
	opt(cfg)
	assert.Equal(t, 5, cfg.NumWorkers)
}
