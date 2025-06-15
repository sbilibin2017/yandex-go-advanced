package repositories

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/storages"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricMemoryGetRepository struct {
	storage *storages.MemoryStorage[types.MetricID, types.Metrics]
}

func NewMetricMemoryGetRepository(
	storage *storages.MemoryStorage[types.MetricID, types.Metrics],
) *MetricMemoryGetRepository {
	return &MetricMemoryGetRepository{
		storage: storage,
	}
}

func (repo *MetricMemoryGetRepository) Get(
	ctx context.Context,
	id types.MetricID,
) (*types.Metrics, error) {
	repo.storage.Mu.RLock()
	defer repo.storage.Mu.RUnlock()

	value, ok := repo.storage.Store[id]
	if !ok {
		return nil, nil
	}
	return &value, nil
}
