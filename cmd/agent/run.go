package main

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/apps"
	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/sbilibin2017/yandex-go-advanced/internal/runners"
)

func run(ctx context.Context, config *configs.AgentConfig) error {
	ctx, stop := runners.NewRunContext(ctx)
	defer stop()

	worker, err := apps.NewAgentApp(config)
	if err != nil {
		return err
	}

	return runners.RunWorker(ctx, worker)

}
