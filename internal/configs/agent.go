package configs

type AgentConfig struct {
	ServerAddress  string
	ServerEndpoint string
	LogLevel       string
	PollInterval   int
	ReportInterval int
	NumWorkers     int
}

type AgentOption func(*AgentConfig)

func NewAgentConfig(opts ...AgentOption) *AgentConfig {
	cfg := &AgentConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
