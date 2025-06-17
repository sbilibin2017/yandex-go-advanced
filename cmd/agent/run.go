package main

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/apps"
	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/sbilibin2017/yandex-go-advanced/internal/logger"
	"github.com/sbilibin2017/yandex-go-advanced/internal/runners"
)

func run(ctx context.Context) error {
	config := configs.NewAgentConfig(
		configs.WithAgentServerAddress(serverAddr),
		configs.WithAgentPollInterval(pollInterval),
		configs.WithAgentReportInterval(reportInterval),
		configs.WithAgentNumWorkers(numWorkers),
		configs.WithAgentLogLevel(logLevel),
	)

	err := logger.Initialize(config.LogLevel)
	if err != nil {
		return err
	}
	logger.Log.Infof("Agent config initialized: %+v", config)

	app, err := apps.NewAgentApp(config)
	if err != nil {
		return err
	}

	return runners.Run(ctx, app)
}
