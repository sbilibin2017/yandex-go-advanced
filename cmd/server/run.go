package main

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/apps"
	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/sbilibin2017/yandex-go-advanced/internal/runners"
)

func run(ctx context.Context, config *configs.ServerConfig) error {
	srv, err := apps.NewServerApp(config)
	if err != nil {
		return err
	}

	ctx, cancel := runners.NewRunContext(ctx)
	defer cancel()

	return runners.RunServer(ctx, srv)
}
