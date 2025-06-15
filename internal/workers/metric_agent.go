package workers

import (
	"context"
	"fmt"
	"math/rand/v2"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

func StartMetricAgentWorker(
	ctx context.Context,
	serverAddress string,
	pollInterval int,
	reportInterval int,
	workerCount int,
) error {
	collectors := []func() []types.Metrics{
		collectRuntimeGaugeMetrics,
		collectRuntimeCounterMetrics,
	}
	metricsCh := pollMetrics(ctx, pollInterval, collectors...)
	errCh := reportMetrics(ctx, serverAddress, sendMetrics, reportInterval, workerCount, metricsCh)
	err := waitForContextOrError(ctx, errCh)
	return err
}

func collectRuntimeGaugeMetrics() []types.Metrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	floatPtr := func(f float64) *float64 { return &f }

	return []types.Metrics{
		{MType: types.Gauge, ID: "Alloc", Value: floatPtr(float64(memStats.Alloc))},
		{MType: types.Gauge, ID: "BuckHashSys", Value: floatPtr(float64(memStats.BuckHashSys))},
		{MType: types.Gauge, ID: "Frees", Value: floatPtr(float64(memStats.Frees))},
		{MType: types.Gauge, ID: "GCCPUFraction", Value: floatPtr(memStats.GCCPUFraction)},
		{MType: types.Gauge, ID: "GCSys", Value: floatPtr(float64(memStats.GCSys))},
		{MType: types.Gauge, ID: "HeapAlloc", Value: floatPtr(float64(memStats.HeapAlloc))},
		{MType: types.Gauge, ID: "HeapIdle", Value: floatPtr(float64(memStats.HeapIdle))},
		{MType: types.Gauge, ID: "HeapInuse", Value: floatPtr(float64(memStats.HeapInuse))},
		{MType: types.Gauge, ID: "HeapObjects", Value: floatPtr(float64(memStats.HeapObjects))},
		{MType: types.Gauge, ID: "HeapReleased", Value: floatPtr(float64(memStats.HeapReleased))},
		{MType: types.Gauge, ID: "HeapSys", Value: floatPtr(float64(memStats.HeapSys))},
		{MType: types.Gauge, ID: "LastGC", Value: floatPtr(float64(memStats.LastGC))},
		{MType: types.Gauge, ID: "Lookups", Value: floatPtr(float64(memStats.Lookups))},
		{MType: types.Gauge, ID: "MCacheInuse", Value: floatPtr(float64(memStats.MCacheInuse))},
		{MType: types.Gauge, ID: "MCacheSys", Value: floatPtr(float64(memStats.MCacheSys))},
		{MType: types.Gauge, ID: "MSpanInuse", Value: floatPtr(float64(memStats.MSpanInuse))},
		{MType: types.Gauge, ID: "MSpanSys", Value: floatPtr(float64(memStats.MSpanSys))},
		{MType: types.Gauge, ID: "Mallocs", Value: floatPtr(float64(memStats.Mallocs))},
		{MType: types.Gauge, ID: "NextGC", Value: floatPtr(float64(memStats.NextGC))},
		{MType: types.Gauge, ID: "NumForcedGC", Value: floatPtr(float64(memStats.NumForcedGC))},
		{MType: types.Gauge, ID: "NumGC", Value: floatPtr(float64(memStats.NumGC))},
		{MType: types.Gauge, ID: "OtherSys", Value: floatPtr(float64(memStats.OtherSys))},
		{MType: types.Gauge, ID: "PauseTotalNs", Value: floatPtr(float64(memStats.PauseTotalNs))},
		{MType: types.Gauge, ID: "StackInuse", Value: floatPtr(float64(memStats.StackInuse))},
		{MType: types.Gauge, ID: "StackSys", Value: floatPtr(float64(memStats.StackSys))},
		{MType: types.Gauge, ID: "Sys", Value: floatPtr(float64(memStats.Sys))},
		{MType: types.Gauge, ID: "TotalAlloc", Value: floatPtr(float64(memStats.TotalAlloc))},
		{MType: types.Gauge, ID: "RandomValue", Value: floatPtr(rand.Float64() * 100)},
	}
}

func collectRuntimeCounterMetrics() []types.Metrics {
	intPtr := func(i int64) *int64 { return &i }

	return []types.Metrics{
		{MType: types.Counter, ID: "PollCount", Delta: intPtr(1)},
	}
}

func pollMetrics(
	ctx context.Context,
	pollInterval int,
	collectors ...func() []types.Metrics,
) <-chan types.Metrics {
	out := make(chan types.Metrics, 100)

	go func() {
		defer close(out)

		interval := time.Duration(pollInterval) * time.Second
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, collect := range collectors {
					metrics := collect()
					for _, metric := range metrics {
						out <- metric
					}
				}
			}
		}
	}()

	return out
}

func sendMetrics(
	ctx context.Context,
	serverAddress string,
	metrics types.Metrics,
) error {
	client := resty.New()

	if !strings.HasPrefix(serverAddress, "http://") && !strings.HasPrefix(serverAddress, "https://") {
		serverAddress = "http://" + serverAddress
	}

	var value string
	switch metrics.MType {
	case types.Counter:
		if metrics.Delta == nil {
			return fmt.Errorf("delta value is nil for Counter metric")
		}
		value = strconv.FormatInt(*metrics.Delta, 10)
	case types.Gauge:
		if metrics.Value == nil {
			return fmt.Errorf("value is nil for Gauge metric")
		}
		value = strconv.FormatFloat(*metrics.Value, 'f', -1, 64)
	default:
		return fmt.Errorf("unknown metric type: %s", metrics.MType)
	}

	url := fmt.Sprintf("%s/update/%s/%s/%s", serverAddress, metrics.MType, metrics.ID, value)

	resp, err := client.R().
		SetContext(ctx).
		Post(url)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to send metrics: %s", resp.Status())
	}

	return nil
}

func reportMetrics(
	ctx context.Context,
	serverAddress string,
	handler func(ctx context.Context, serverAddress string, metrics types.Metrics) error,
	reportInterval int,
	workerCount int,
	in <-chan types.Metrics,
) <-chan error {
	errCh := make(chan error, 100)
	jobs := make(chan types.Metrics, 100)

	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case metric, ok := <-jobs:
				if !ok {
					return
				}
				if err := handler(ctx, serverAddress, metric); err != nil {
					errCh <- err
				}
			}
		}
	}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	go func() {
		defer close(jobs)
		ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
		defer ticker.Stop()

		var buffer []types.Metrics

		flush := func() {
			for _, metric := range buffer {
				jobs <- metric
			}
			buffer = buffer[:0]
		}

		for {
			select {
			case <-ctx.Done():
				flush()
				return
			case metric, ok := <-in:
				if !ok {
					flush()
					return
				}
				buffer = append(buffer, metric)
			case <-ticker.C:
				flush()
			}
		}
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	return errCh
}

func waitForContextOrError(ctx context.Context, errCh <-chan error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err, ok := <-errCh:
		if !ok {
			return nil
		}
		return err
	}
}
