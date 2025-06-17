// Package configs provides configuration structures and constructors
// for the agent application.
package configs

// AgentConfig holds configuration parameters for the agent.
type AgentConfig struct {
	ServerAddress  string // Address of the server to send metrics to
	LogLevel       string // Logging level (e.g., debug, info, warn, error)
	PollInterval   int    // Time interval (in seconds) between metric polling
	ReportInterval int    // Time interval (in seconds) between sending metrics
	NumWorkers     int    // Number of concurrent workers for sending metrics
}

// AgentOption defines a function that modifies an AgentConfig.
// It is used for functional-style configuration.
type AgentOption func(*AgentConfig)

// NewAgentConfig creates a new AgentConfig and applies any number of
// AgentOption functions to it.
//
// Example usage:
//
//	cfg := NewAgentConfig(
//	    func(c *AgentConfig) { c.ServerAddress = "localhost:8080" },
//	    func(c *AgentConfig) { c.PollInterval = 10 },
//	)
//
// Parameters:
//   - opts: Variadic list of AgentOption functions to configure the agent.
//
// Returns:
//   - A pointer to the fully constructed AgentConfig.
func NewAgentConfig(opts ...AgentOption) *AgentConfig {
	cfg := &AgentConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithServerAddress sets the ServerAddress field.
func WithAgentServerAddress(addr string) AgentOption {
	return func(cfg *AgentConfig) {
		cfg.ServerAddress = addr
	}
}

// WithLogLevel sets the LogLevel field.
func WithAgentLogLevel(level string) AgentOption {
	return func(cfg *AgentConfig) {
		cfg.LogLevel = level
	}
}

// WithPollInterval sets the PollInterval field.
func WithAgentPollInterval(interval int) AgentOption {
	return func(cfg *AgentConfig) {
		cfg.PollInterval = interval
	}
}

// WithReportInterval sets the ReportInterval field.
func WithAgentReportInterval(interval int) AgentOption {
	return func(cfg *AgentConfig) {
		cfg.ReportInterval = interval
	}
}

// WithNumWorkers sets the NumWorkers field.
func WithAgentNumWorkers(workers int) AgentOption {
	return func(cfg *AgentConfig) {
		cfg.NumWorkers = workers
	}
}
