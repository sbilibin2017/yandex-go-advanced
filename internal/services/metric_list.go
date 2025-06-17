package services

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricListLister defines an interface for listing metrics.
type MetricListLister interface {
	// List returns a slice of all available metrics or an error if something goes wrong.
	List(ctx context.Context) ([]types.Metrics, error)
}

// MetricListService provides functionality to retrieve a list of metrics.
type MetricListService struct {
	lister MetricListLister
}

// NewMetricListService creates a new MetricListService using the provided MetricListLister.
func NewMetricListService(
	lister MetricListLister,
) *MetricListService {
	return &MetricListService{lister: lister}
}

// List fetches the list of metrics by delegating to the underlying MetricListLister.
// It returns the slice of metrics or an error if the operation fails.
func (svc *MetricListService) List(
	ctx context.Context,
) ([]types.Metrics, error) {
	metrics, err := svc.lister.List(ctx)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
