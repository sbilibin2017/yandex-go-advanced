package facades

import (
	"context"
	"fmt"
	"strconv"
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

func (m *MetricUpdateFacade) Update(ctx context.Context, metrics []types.Metrics) error {
	serverAddress := m.serverAddress
	if !strings.HasPrefix(serverAddress, "http://") && !strings.HasPrefix(serverAddress, "https://") {
		serverAddress = "http://" + serverAddress
	}

	for _, metric := range metrics {
		var value string
		switch metric.Type {
		case types.Counter:
			if metric.Delta == nil {
				return fmt.Errorf("delta value is nil for Counter metric")
			}
			value = strconv.FormatInt(*metric.Delta, 10)
		case types.Gauge:
			if metric.Value == nil {
				return fmt.Errorf("value is nil for Gauge metric")
			}
			value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
		default:
			return fmt.Errorf("unknown metric type: %s", metric.Type)
		}

		url := fmt.Sprintf("%s/update/%s/%s/%s", serverAddress, metric.Type, metric.ID, value)

		resp, err := m.client.R().
			SetContext(ctx).
			Post(url)
		if err != nil {
			return fmt.Errorf("failed to send metric %s/%s: %w", metric.Type, metric.ID, err)
		}

		if resp.IsError() {
			return fmt.Errorf("failed to send metric %s/%s: %s", metric.Type, metric.ID, resp.Status())
		}
	}

	return nil
}
