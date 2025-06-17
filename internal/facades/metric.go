package facades

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sbilibin2017/yandex-go-advanced/internal/logger"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricUpdateFacade provides a simplified interface for sending
// metric update requests to a remote server.
type MetricUpdateFacade struct {
	serverAddress string
	client        *resty.Client
}

// NewMetricUpdateFacade creates and returns a new MetricUpdateFacade.
// It initializes an HTTP client and accepts the server address to which
// the metrics will be sent.
func NewMetricUpdateFacade(serverAddress string) *MetricUpdateFacade {
	client := resty.New()
	return &MetricUpdateFacade{
		serverAddress: serverAddress,
		client:        client,
	}
}

// Update sends the provided slice of metrics to the configured server address
// using individual POST requests to the /update/ endpoint.
//
// It ensures the server address has the proper protocol prefix,
// and returns an error if any request fails or if the response indicates a failure.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout.
//   - metrics: Slice of metric pointers to be sent.
//
// Returns:
//   - An error if any update fails or the server responds with an error status.
func (m *MetricUpdateFacade) Update(ctx context.Context, metrics []*types.Metrics) error {
	serverAddress := m.serverAddress
	if !strings.HasPrefix(serverAddress, "http://") && !strings.HasPrefix(serverAddress, "https://") {
		serverAddress = "http://" + serverAddress
	}

	url := fmt.Sprintf("%s/update/", serverAddress)

	for _, metric := range metrics {
		body, err := compressMetrics(metric)
		if err != nil {
			return err
		}
		resp, err := m.client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(body).
			Post(url)

		if err != nil {
			logger.Log.Errorf("Failed to send metrics update request for metric ID=%s: %v", metric.ID, err)
			return fmt.Errorf("failed to send metrics update request: %w", err)
		}

		if resp.IsError() {
			logger.Log.Errorf("Metrics update request failed for metric ID=%s: %s", metric.ID, resp.Status())
			return fmt.Errorf("metrics update request failed: %s", resp.Status())
		}

	}

	return nil
}

func compressMetrics(metrics *types.Metrics) ([]byte, error) {
	data, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	defer func() {
		_ = gz.Close()
	}()

	_, err = gz.Write(data)
	if err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
