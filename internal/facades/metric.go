package facades

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricUpdateFacade struct {
	serverAddress string
	client        *resty.Client
}

func NewMetricUpdateFacade(serverAddress string) *MetricUpdateFacade {
	client := resty.New()
	return &MetricUpdateFacade{
		serverAddress: serverAddress,
		client:        client,
	}
}

func (m *MetricUpdateFacade) Update(ctx context.Context, metrics []*types.Metrics) error {
	serverAddress := m.serverAddress
	if !strings.HasPrefix(serverAddress, "http://") && !strings.HasPrefix(serverAddress, "https://") {
		serverAddress = "http://" + serverAddress
	}

	url := fmt.Sprintf("%s/update/", serverAddress)

	for _, metric := range metrics {
		resp, err := m.client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetBody(metric).
			Post(url)

		if err != nil {
			return fmt.Errorf("failed to send metrics update request: %w", err)
		}

		if resp.IsError() {
			return fmt.Errorf("metrics update request failed: %s", resp.Status())
		}

	}

	return nil
}
