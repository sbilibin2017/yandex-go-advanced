package main

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/apps"
	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/sbilibin2017/yandex-go-advanced/internal/logger"
	"github.com/sbilibin2017/yandex-go-advanced/internal/runners"
)

func run(ctx context.Context, config *configs.AgentConfig) error {
	err := logger.Initialize(config.LogLevel)
	if err != nil {
		return err
	}

	worker, err := apps.NewAgentApp(config)
	if err != nil {
		return err
	}

	runners.RunWorker(ctx, worker)

	return nil

}
