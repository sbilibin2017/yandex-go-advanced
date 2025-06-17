// Package configs provides configuration structures and constructors
// for the server application.
package configs

// ServerConfig holds configuration parameters for the HTTP server.
type ServerConfig struct {
	Address  string // Address on which the server listens (e.g., ":8080")
	LogLevel string // Logging level (e.g., debug, info, warn, error)
}

// ServerOption defines a function that modifies a ServerConfig.
// It is used for functional-style configuration.
type ServerOption func(*ServerConfig)

// NewServerConfig creates a new ServerConfig and applies any number of
// ServerOption functions to it.
//
// Example usage:
//
//	cfg := NewServerConfig(
//	    func(c *ServerConfig) { c.Address = ":8080" },
//	    func(c *ServerConfig) { c.LogLevel = "info" },
//	)
//
// Parameters:
//   - opts: Variadic list of ServerOption functions to configure the server.
//
// Returns:
//   - A pointer to the fully constructed ServerConfig.
func NewServerConfig(opts ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithAddress sets the server listening address.
func WithServerAddress(addr string) ServerOption {
	return func(c *ServerConfig) {
		c.Address = addr
	}
}

// WithLogLevel sets the logging level.
func WithServerLogLevel(level string) ServerOption {
	return func(c *ServerConfig) {
		c.LogLevel = level
	}
}
