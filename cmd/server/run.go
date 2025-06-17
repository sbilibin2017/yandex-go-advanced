package main

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/apps"
	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/sbilibin2017/yandex-go-advanced/internal/logger"
	"github.com/sbilibin2017/yandex-go-advanced/internal/runners"
)

func run(ctx context.Context) error {
	config := configs.NewServerConfig(
		configs.WithServerAddress(addr),
		configs.WithServerLogLevel(logLevel),
	)

	err := logger.Initialize(config.LogLevel)
	if err != nil {
		logger.Log.Errorf("Failed to initialize logger: %v", err)
		return err
	}
	logger.Log.Infof("Server config initialized: %+v", config)

	app, err := apps.NewServerApp(config)
	if err != nil {
		logger.Log.Errorf("Failed to create server app: %v", err)
		return err
	}

	err = runners.Run(ctx, app)
	if err != nil {
		logger.Log.Errorf("Error running the app: %v", err)
		return err
	}

	return nil
}
