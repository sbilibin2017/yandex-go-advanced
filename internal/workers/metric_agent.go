package workers

import (
	"context"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricUpdater interface {
	Update(ctx context.Context, metrics []types.Metrics) error
}

func NewMetricAgentWorker(
	updater MetricUpdater,
	pollInterval int,
	reportInterval int,
	workerCount int,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		return startMetricAgentWorker(ctx, updater, pollInterval, reportInterval, workerCount)
	}
}

func startMetricAgentWorker(
	ctx context.Context,
	updater MetricUpdater,
	pollInterval int,
	reportInterval int,
	workerCount int,
) error {
	collectors := []func() []types.Metrics{
		collectRuntimeGaugeMetrics,
		collectRuntimeCounterMetrics,
	}
	metricsCh := pollMetrics(ctx, pollInterval, collectors...)
	errCh := reportMetrics(ctx, updater, reportInterval, workerCount, metricsCh)
	err := waitForContextOrError(ctx, errCh)
	return err
}

func collectRuntimeGaugeMetrics() []types.Metrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	floatPtr := func(f float64) *float64 { return &f }

	return []types.Metrics{
		{Type: types.Gauge, ID: "Alloc", Value: floatPtr(float64(memStats.Alloc))},
		{Type: types.Gauge, ID: "BuckHashSys", Value: floatPtr(float64(memStats.BuckHashSys))},
		{Type: types.Gauge, ID: "Frees", Value: floatPtr(float64(memStats.Frees))},
		{Type: types.Gauge, ID: "GCCPUFraction", Value: floatPtr(memStats.GCCPUFraction)},
		{Type: types.Gauge, ID: "GCSys", Value: floatPtr(float64(memStats.GCSys))},
		{Type: types.Gauge, ID: "HeapAlloc", Value: floatPtr(float64(memStats.HeapAlloc))},
		{Type: types.Gauge, ID: "HeapIdle", Value: floatPtr(float64(memStats.HeapIdle))},
		{Type: types.Gauge, ID: "HeapInuse", Value: floatPtr(float64(memStats.HeapInuse))},
		{Type: types.Gauge, ID: "HeapObjects", Value: floatPtr(float64(memStats.HeapObjects))},
		{Type: types.Gauge, ID: "HeapReleased", Value: floatPtr(float64(memStats.HeapReleased))},
		{Type: types.Gauge, ID: "HeapSys", Value: floatPtr(float64(memStats.HeapSys))},
		{Type: types.Gauge, ID: "LastGC", Value: floatPtr(float64(memStats.LastGC))},
		{Type: types.Gauge, ID: "Lookups", Value: floatPtr(float64(memStats.Lookups))},
		{Type: types.Gauge, ID: "MCacheInuse", Value: floatPtr(float64(memStats.MCacheInuse))},
		{Type: types.Gauge, ID: "MCacheSys", Value: floatPtr(float64(memStats.MCacheSys))},
		{Type: types.Gauge, ID: "MSpanInuse", Value: floatPtr(float64(memStats.MSpanInuse))},
		{Type: types.Gauge, ID: "MSpanSys", Value: floatPtr(float64(memStats.MSpanSys))},
		{Type: types.Gauge, ID: "Mallocs", Value: floatPtr(float64(memStats.Mallocs))},
		{Type: types.Gauge, ID: "NextGC", Value: floatPtr(float64(memStats.NextGC))},
		{Type: types.Gauge, ID: "NumForcedGC", Value: floatPtr(float64(memStats.NumForcedGC))},
		{Type: types.Gauge, ID: "NumGC", Value: floatPtr(float64(memStats.NumGC))},
		{Type: types.Gauge, ID: "OtherSys", Value: floatPtr(float64(memStats.OtherSys))},
		{Type: types.Gauge, ID: "PauseTotalNs", Value: floatPtr(float64(memStats.PauseTotalNs))},
		{Type: types.Gauge, ID: "StackInuse", Value: floatPtr(float64(memStats.StackInuse))},
		{Type: types.Gauge, ID: "StackSys", Value: floatPtr(float64(memStats.StackSys))},
		{Type: types.Gauge, ID: "Sys", Value: floatPtr(float64(memStats.Sys))},
		{Type: types.Gauge, ID: "TotalAlloc", Value: floatPtr(float64(memStats.TotalAlloc))},
		{Type: types.Gauge, ID: "RandomValue", Value: floatPtr(rand.Float64() * 100)},
	}
}

func collectRuntimeCounterMetrics() []types.Metrics {
	intPtr := func(i int64) *int64 { return &i }

	return []types.Metrics{
		{Type: types.Counter, ID: "PollCount", Delta: intPtr(1)},
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

func reportMetrics(
	ctx context.Context,
	updater MetricUpdater,
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
				if err := updater.Update(ctx, []types.Metrics{metric}); err != nil {
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
