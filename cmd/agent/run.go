package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sbilibin2017/yandex-go-advanced/internal/workers"
)

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	return workers.StartMetricAgentWorker(
		ctx,
		addressFlag,
		pollIntervalFlag,
		reportIntervalFlag,
		numWorkersFlag,
	)

}
