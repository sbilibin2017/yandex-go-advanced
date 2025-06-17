// Package configs provides configuration structures and constructors
// for the server application.
package configs

// ServerConfig holds configuration parameters for the HTTP server
// and metrics storage.
type ServerConfig struct {
	Address         string // Address on which the server listens (e.g., ":8080")
	LogLevel        string // Logging level (e.g., debug, info, warn, error)
	StoreInterval   int    // Interval in seconds to save metrics to disk (0 = sync save)
	FileStoragePath string // File path to store metrics
	Restore         bool   // Whether to restore metrics from file on startup
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
//	    WithServerAddress(":8080"),
//	    WithServerLogLevel("info"),
//	    WithStoreInterval(300),
//	    WithFileStoragePath("metrics.json"),
//	    WithRestore(true),
//	)
//
// Parameters:
//   - opts: Variadic list of ServerOption functions to configure the server.
//
// Returns:
//   - A pointer to the fully constructed ServerConfig.
func NewServerConfig(opts ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{
		StoreInterval:   300,            // default 300 seconds
		FileStoragePath: "metrics.json", // default filename
		Restore:         false,          // default false
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithServerAddress sets the server listening address.
func WithServerAddress(addr string) ServerOption {
	return func(c *ServerConfig) {
		c.Address = addr
	}
}

// WithServerLogLevel sets the logging level.
func WithServerLogLevel(level string) ServerOption {
	return func(c *ServerConfig) {
		c.LogLevel = level
	}
}

// WithStoreInterval sets the interval in seconds for storing metrics.
func WithStoreInterval(interval int) ServerOption {
	return func(c *ServerConfig) {
		c.StoreInterval = interval
	}
}

// WithFileStoragePath sets the file path where metrics are stored.
func WithFileStoragePath(path string) ServerOption {
	return func(c *ServerConfig) {
		c.FileStoragePath = path
	}
}

// WithRestore sets whether to restore metrics from file on startup.
func WithRestore(restore bool) ServerOption {
	return func(c *ServerConfig) {
		c.Restore = restore
	}
}
