package main

import (
	"flag"
	"os"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
)

func parseFlags() (*configs.ServerConfig, error) {
	fs := flag.NewFlagSet("server", flag.ExitOnError)

	options := []configs.ServerOption{
		withAddr(fs),
		withLogLevel(fs),
	}

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	return configs.NewServerConfig(options...), nil
}

func withAddr(fs *flag.FlagSet) configs.ServerOption {
	var addrFlag string
	fs.StringVar(&addrFlag, "a", "localhost:8080", "address and port to run server")

	return func(cfg *configs.ServerConfig) {
		if env := os.Getenv("ADDRESS"); env != "" {
			cfg.Address = env
			return
		}
		cfg.Address = addrFlag
	}
}

func withLogLevel(fs *flag.FlagSet) configs.ServerOption {
	var levelFlag string
	fs.StringVar(&levelFlag, "l", "info", "log level")

	return func(cfg *configs.ServerConfig) {
		if env := os.Getenv("LOG_LEVEL"); env != "" {
			cfg.LogLevel = env
			return
		}
		cfg.LogLevel = levelFlag
	}
}
