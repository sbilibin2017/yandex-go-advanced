package repositories

import (
	"context"
	"sort"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricMemoryListRepository struct{}

func NewMetricMemoryListRepository() *MetricMemoryListRepository {
	return &MetricMemoryListRepository{}
}

func (repo *MetricMemoryListRepository) List(ctx context.Context) ([]types.Metrics, error) {
	mu.RLock()
	defer mu.RUnlock()

	list := make([]types.Metrics, 0, len(metrics))
	for _, m := range metrics {
		list = append(list, m)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})

	return list, nil
}
