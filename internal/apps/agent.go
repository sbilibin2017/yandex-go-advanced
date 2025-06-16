package apps

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/sbilibin2017/yandex-go-advanced/internal/facades"
	"github.com/sbilibin2017/yandex-go-advanced/internal/workers"
)

func NewAgentApp(
	config *configs.AgentConfig,
) (func(ctx context.Context) error, error) {
	metricUpdateFacade := facades.NewMetricUpdateFacade(config.ServerAddress)

	worker := workers.NewMetricAgentWorker(
		metricUpdateFacade,
		config.PollInterval,
		config.ReportInterval,
		config.NumWorkers,
	)

	return worker, nil
}
