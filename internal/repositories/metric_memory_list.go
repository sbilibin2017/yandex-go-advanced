package repositories

import (
	"context"
	"sort"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricMemoryListRepository provides in-memory listing of all stored metrics.
type MetricMemoryListRepository struct{}

// NewMetricMemoryListRepository creates and returns a new MetricMemoryListRepository instance.
func NewMetricMemoryListRepository() *MetricMemoryListRepository {
	return &MetricMemoryListRepository{}
}

// List returns all metrics currently stored in memory, sorted by their MetricID.
//
// Parameters:
//   - ctx: Context for cancellation and deadlines (not used in current implementation).
//
// Returns:
//   - A slice of Metrics sorted by their ID.
//   - An error if the operation fails (currently always nil as no error handling is implemented).
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
