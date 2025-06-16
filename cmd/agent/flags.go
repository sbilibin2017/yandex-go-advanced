package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
)

func parseFlags() (*configs.AgentConfig, error) {
	fs := flag.NewFlagSet("agent", flag.ExitOnError)

	options := []configs.AgentOption{
		withServerAddress(fs),
		withPollInterval(fs),
		withReportInterval(fs),
		withNumWorkers(fs),
		withLogLevel(fs),
	}

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	return configs.NewAgentConfig(options...), nil
}

func withServerAddress(fs *flag.FlagSet) configs.AgentOption {
	var addrFlag string
	fs.StringVar(&addrFlag, "a", "localhost:8080", "server address")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("ADDRESS"); env != "" {
			cfg.ServerAddress = env
			return
		}
		cfg.ServerAddress = addrFlag
	}
}

func withPollInterval(fs *flag.FlagSet) configs.AgentOption {
	var pollFlag int
	fs.IntVar(&pollFlag, "p", 2, "polling interval in seconds")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("POLL_INTERVAL"); env != "" {
			if v, err := strconv.Atoi(env); err == nil {
				cfg.PollInterval = v
				return
			}
		}
		cfg.PollInterval = pollFlag
	}
}

func withReportInterval(fs *flag.FlagSet) configs.AgentOption {
	var reportFlag int
	fs.IntVar(&reportFlag, "r", 10, "reporting interval in seconds")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("REPORT_INTERVAL"); env != "" {
			if v, err := strconv.Atoi(env); err == nil {
				cfg.ReportInterval = v
				return
			}
		}
		cfg.ReportInterval = reportFlag
	}
}

func withNumWorkers(fs *flag.FlagSet) configs.AgentOption {
	var workersFlag int
	fs.IntVar(&workersFlag, "workers", 4, "number of workers")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("NUM_WORKERS"); env != "" {
			if v, err := strconv.Atoi(env); err == nil {
				cfg.NumWorkers = v
				return
			}
		}
		cfg.NumWorkers = workersFlag
	}
}

func withLogLevel(fs *flag.FlagSet) configs.AgentOption {
	var levelFlag string
	fs.StringVar(&levelFlag, "l", "info", "log level")

	return func(cfg *configs.AgentConfig) {
		if env := os.Getenv("LOG_LEVEL"); env != "" {
			cfg.LogLevel = env
			return
		}
		cfg.LogLevel = levelFlag
	}
}
