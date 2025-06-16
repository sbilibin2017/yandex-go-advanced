package configs

type ServerConfig struct {
	Address  string
	LogLevel string
}

type ServerOption func(*ServerConfig)

func NewServerConfig(opts ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
