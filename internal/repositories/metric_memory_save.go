package repositories

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricMemorySaveRepository struct{}

func NewMetricMemorySaveRepository() *MetricMemorySaveRepository {
	return &MetricMemorySaveRepository{}
}

func (repo *MetricMemorySaveRepository) Save(
	ctx context.Context,
	m types.Metrics,
) error {
	mu.Lock()
	defer mu.Unlock()

	metrics[types.MetricID{ID: m.ID, MType: m.MType}] = m
	return nil
}
