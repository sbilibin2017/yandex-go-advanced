package services

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricGetGetter defines the interface for retrieving a metric by its ID.
type MetricGetGetter interface {
	// Get fetches the metric identified by the given MetricID.
	// It returns the metric if found, or an error otherwise.
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

// MetricGetService provides metric retrieval services by delegating to a MetricGetGetter.
type MetricGetService struct {
	getter MetricGetGetter
}

// NewMetricGetService creates a new MetricGetService with the provided MetricGetGetter.
func NewMetricGetService(
	getter MetricGetGetter,
) *MetricGetService {
	return &MetricGetService{getter: getter}
}

// Get retrieves a metric by its MetricID using the underlying MetricGetGetter.
func (svc *MetricGetService) Get(
	ctx context.Context,
	id types.MetricID,
) (*types.Metrics, error) {
	metric, err := svc.getter.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return metric, nil
}
