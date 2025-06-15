package repositories

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricMemoryGetRepository struct{}

func NewMetricMemoryGetRepository() *MetricMemoryGetRepository {
	return &MetricMemoryGetRepository{}
}

func (repo *MetricMemoryGetRepository) Get(
	ctx context.Context,
	id types.MetricID,
) (*types.Metrics, error) {
	mu.RLock()
	defer mu.RUnlock()

	value, ok := metrics[id]
	if !ok {
		return nil, nil
	}
	return &value, nil
}
