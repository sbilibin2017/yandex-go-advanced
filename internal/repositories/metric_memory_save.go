package repositories

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/storages"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricMemorySaveRepository struct {
	storage *storages.MemoryStorage[types.MetricID, types.Metrics]
}

func NewMetricMemorySaveRepository(
	storage *storages.MemoryStorage[types.MetricID, types.Metrics],
) *MetricMemorySaveRepository {
	return &MetricMemorySaveRepository{
		storage: storage,
	}
}

func (repo *MetricMemorySaveRepository) Save(
	ctx context.Context,
	metrics types.Metrics,
) error {
	repo.storage.Mu.Lock()
	defer repo.storage.Mu.Unlock()

	repo.storage.Store[types.MetricID{ID: metrics.ID, MType: metrics.MType}] = metrics

	return nil
}
