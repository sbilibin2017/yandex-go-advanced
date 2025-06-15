package repositories

import (
	"context"
	"sort"

	"github.com/sbilibin2017/yandex-go-advanced/internal/storages"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricMemoryListRepository struct {
	storage *storages.MemoryStorage[types.MetricID, types.Metrics]
}

func NewMetricMemoryListRepository(
	storage *storages.MemoryStorage[types.MetricID, types.Metrics],
) *MetricMemoryListRepository {
	return &MetricMemoryListRepository{
		storage: storage,
	}
}

func (repo *MetricMemoryListRepository) List(ctx context.Context) ([]types.Metrics, error) {
	repo.storage.Mu.RLock()
	defer repo.storage.Mu.RUnlock()

	metrics := make([]types.Metrics, 0, len(repo.storage.Store))
	for _, metric := range repo.storage.Store {
		metrics = append(metrics, metric)
	}

	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].ID < metrics[j].ID
	})

	return metrics, nil
}
