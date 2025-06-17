package workers

import (
	"context"
	"time"

	"github.com/sbilibin2017/yandex-go-advanced/internal/logger"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricContextSaver defines an interface for saving metrics to the in-memory or context storage.
type MetricContextSaver interface {
	// Save persists the given metric data into the context storage.
	Save(ctx context.Context, metrics types.Metrics) error
}

// MetricContextLister defines an interface for listing metrics from the in-memory or context storage.
type MetricContextLister interface {
	// List returns all metrics currently stored in the context.
	List(ctx context.Context) ([]types.Metrics, error)
}

// MetricFileSaver defines an interface for saving metrics to a file storage.
type MetricFileSaver interface {
	// Save persists the given metric data into a file.
	Save(ctx context.Context, metrics types.Metrics) error
}

// MetricFileLister defines an interface for listing metrics from a file storage.
type MetricFileLister interface {
	// List returns all metrics stored in the file.
	List(ctx context.Context) ([]types.Metrics, error)
}

// NewMetricServerWorker returns a function that runs the metric server worker with the given parameters.
// The worker periodically saves metrics to file and optionally restores metrics from file on startup.
//
// saverContext: storage interface to save metrics in memory/context
// listerContext: storage interface to list metrics from memory/context
// saverFile: storage interface to save metrics to file
// listerFile: storage interface to list metrics from file
// storeInterval: interval in seconds for periodic saving of metrics; 0 disables periodic saving (only saves on shutdown)
// restore: if true, loads metrics from file on startup before starting the periodic saving loop
func NewMetricServerWorker(
	saverContext MetricContextSaver,
	listerContext MetricContextLister,
	saverFile MetricFileSaver,
	listerFile MetricFileLister,
	storeInterval int,
	restore bool,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		startMetricServerWorker(
			ctx,
			saverContext,
			listerContext,
			saverFile,
			listerFile,
			storeInterval,
			restore,
		)
	}
}

// startMetricServerWorker runs the main worker loop to periodically save metrics to file and optionally restore them on startup.
//
// The worker listens for context cancellation to gracefully shutdown and save metrics one last time.
//
// If storeInterval is 0, the worker only saves metrics once on shutdown.
func startMetricServerWorker(
	ctx context.Context,
	saverContext MetricContextSaver,
	listerContext MetricContextLister,
	saverFile MetricFileSaver,
	listerFile MetricFileLister,
	storeInterval int,
	restore bool,
) {
	if restore {
		if err := loadFromFile(ctx, saverContext, listerFile); err != nil {
			logger.Log.Errorw("failed to restore metrics from file", "error", err)
		}
	}

	if storeInterval == 0 {
		<-ctx.Done()
		if err := saveToFile(ctx, listerContext, saverFile); err != nil {
			logger.Log.Errorw("failed to save metrics to file during shutdown", "error", err)
		}
		return
	}

	ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if err := saveToFile(ctx, listerContext, saverFile); err != nil {
				logger.Log.Errorw("failed to save metrics to file during shutdown", "error", err)
			}
			return
		case <-ticker.C:
			if err := saveToFile(ctx, listerContext, saverFile); err != nil {
				logger.Log.Errorw("failed to save metrics to file", "error", err)
			}
		}
	}
}

// loadFromFile loads metrics from the file using the listerFile interface and saves them into the context storage.
//
// This is typically called during startup if restore is enabled.
func loadFromFile(ctx context.Context, saverContext MetricContextSaver, listerFile MetricFileLister) error {
	metrics, err := listerFile.List(ctx)
	if err != nil {
		return err
	}

	for _, m := range metrics {
		if err := saverContext.Save(ctx, m); err != nil {
			logger.Log.Errorw("failed to save restored metric to context storage", "error", err)
		}
	}
	return nil
}

// saveToFile retrieves all metrics from the context storage and saves them into the file using the saverFile interface.
func saveToFile(ctx context.Context, listerContext MetricContextLister, saverFile MetricFileSaver) error {
	metrics, err := listerContext.List(ctx)
	if err != nil {
		return err
	}

	for _, m := range metrics {
		if err := saverFile.Save(ctx, m); err != nil {
			logger.Log.Errorw("failed to save metric to file", "error", err)
		}
	}
	return nil
}
