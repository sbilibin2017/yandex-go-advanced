package services

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricUpdateSaver defines an interface for saving metrics.
type MetricUpdateSaver interface {
	// Save persists the given metric data.
	Save(ctx context.Context, metrics types.Metrics) error
}

// MetricUpdateGetter defines an interface for retrieving metrics by ID.
type MetricUpdateGetter interface {
	// Get retrieves a metric by its ID. Returns nil if not found.
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

// MetricUpdateService provides methods for updating metrics,
// combining retrieving and saving functionality.
type MetricUpdateService struct {
	saver  MetricUpdateSaver
	getter MetricUpdateGetter
}

// NewMetricUpdateService creates a new MetricUpdateService with the provided saver and getter.
func NewMetricUpdateService(
	saver MetricUpdateSaver,
	getter MetricUpdateGetter,
) *MetricUpdateService {
	return &MetricUpdateService{saver: saver, getter: getter}
}

// Update processes and saves a slice of metrics.
// For counter-type metrics, it sums the existing delta with the new one before saving.
// Returns the updated slice of metrics or an error.
func (svc *MetricUpdateService) Update(
	ctx context.Context,
	metrics []*types.Metrics,
) ([]*types.Metrics, error) {
	for idx, m := range metrics {
		switch m.Type {
		case types.Counter:
			if err := updateCounterMetric(ctx, svc.getter, m); err != nil {
				return nil, err
			}
		}

		err := svc.saver.Save(ctx, *m)
		if err != nil {
			return nil, err
		}
		metrics[idx] = m
	}

	return metrics, nil
}

// updateCounterMetric retrieves the existing counter metric and sums its delta value with the incoming one.
func updateCounterMetric(
	ctx context.Context,
	getter MetricUpdateGetter,
	metric *types.Metrics,
) error {
	existing, err := getter.Get(ctx, types.MetricID{ID: metric.ID, Type: metric.Type})
	if err != nil {
		return err
	}

	if existing != nil && existing.Delta != nil && metric.Delta != nil {
		*metric.Delta += *existing.Delta
	}

	return nil
}
