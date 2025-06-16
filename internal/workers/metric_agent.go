package workers

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"github.com/sbilibin2017/yandex-go-advanced/internal/logger"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricUpdater interface {
	Update(ctx context.Context, metrics []*types.Metrics) error
}

func NewMetricAgentWorker(
	updater MetricUpdater,
	pollInterval int,
	reportInterval int,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		startMetricAgentWorker(ctx, updater, pollInterval, reportInterval)
	}
}

func startMetricAgentWorker(
	ctx context.Context,
	updater MetricUpdater,
	pollInterval int,
	reportInterval int,
) {
	pollCh := collectRuntimeMetrics(ctx, pollInterval)
	reportCh := updateMetrics(ctx, reportInterval, updater, pollCh)
	logErrors(ctx, reportCh)
}

func collectRuntimeMetrics(ctx context.Context, pollInterval int) <-chan *types.Metrics {
	out := make(chan *types.Metrics)

	go func() {
		defer close(out)
		ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ms := &runtime.MemStats{}
				runtime.ReadMemStats(ms)

				sendGauge := func(name string, val float64) {
					out <- &types.Metrics{ID: name, Type: "gauge", Value: &val}
				}

				sendGauge("Alloc", float64(ms.Alloc))
				sendGauge("BuckHashSys", float64(ms.BuckHashSys))
				sendGauge("Frees", float64(ms.Frees))
				sendGauge("GCCPUFraction", ms.GCCPUFraction)
				sendGauge("GCSys", float64(ms.GCSys))
				sendGauge("HeapAlloc", float64(ms.HeapAlloc))
				sendGauge("HeapIdle", float64(ms.HeapIdle))
				sendGauge("HeapInuse", float64(ms.HeapInuse))
				sendGauge("HeapObjects", float64(ms.HeapObjects))
				sendGauge("HeapReleased", float64(ms.HeapReleased))
				sendGauge("HeapSys", float64(ms.HeapSys))
				sendGauge("LastGC", float64(ms.LastGC))
				sendGauge("Lookups", float64(ms.Lookups))
				sendGauge("MCacheInuse", float64(ms.MCacheInuse))
				sendGauge("MCacheSys", float64(ms.MCacheSys))
				sendGauge("MSpanInuse", float64(ms.MSpanInuse))
				sendGauge("MSpanSys", float64(ms.MSpanSys))
				sendGauge("Mallocs", float64(ms.Mallocs))
				sendGauge("NextGC", float64(ms.NextGC))
				sendGauge("NumForcedGC", float64(ms.NumForcedGC))
				sendGauge("NumGC", float64(ms.NumGC))
				sendGauge("OtherSys", float64(ms.OtherSys))
				sendGauge("PauseTotalNs", float64(ms.PauseTotalNs))
				sendGauge("StackInuse", float64(ms.StackInuse))
				sendGauge("StackSys", float64(ms.StackSys))
				sendGauge("Sys", float64(ms.Sys))
				sendGauge("TotalAlloc", float64(ms.TotalAlloc))

				c := int64(1)
				out <- &types.Metrics{ID: "PollCount", Type: "counter", Delta: &c}

				rv := rand.Float64()
				out <- &types.Metrics{ID: "RandomValue", Type: "gauge", Value: &rv}
			}
		}
	}()

	return out
}

func updateMetrics(
	ctx context.Context,
	reportInterval int,
	updater MetricUpdater,
	in <-chan *types.Metrics,
) <-chan error {
	errCh := make(chan error)
	ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)

	go func() {
		defer close(errCh)
		defer ticker.Stop()

		var buffer []*types.Metrics

		for {
			select {
			case <-ctx.Done():
				if len(buffer) > 0 {
					if err := updater.Update(ctx, buffer); err != nil {
						errCh <- err
					}
				}
				return

			case m, ok := <-in:
				if !ok {
					if len(buffer) > 0 {
						if err := updater.Update(ctx, buffer); err != nil {
							errCh <- err
						}
					}
					return
				}
				buffer = append(buffer, m)

			case <-ticker.C:
				if len(buffer) > 0 {
					if err := updater.Update(ctx, buffer); err != nil {
						errCh <- err
					}
					buffer = buffer[:0]
				}
			}
		}
	}()

	return errCh
}

func logErrors(ctx context.Context, errCh <-chan error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err, ok := <-errCh:
				if !ok {
					return
				}
				if err != nil {
					logger.Log.Error("update error: ", err)
				}
			}
		}
	}()
}
