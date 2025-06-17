package repositories

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricMemoryGetRepository provides in-memory retrieval of metrics.
type MetricMemoryGetRepository struct{}

// NewMetricMemoryGetRepository creates and returns a new MetricMemoryGetRepository instance.
func NewMetricMemoryGetRepository() *MetricMemoryGetRepository {
	return &MetricMemoryGetRepository{}
}

// Get retrieves a metric by its ID from the in-memory storage.
//
// Parameters:
//   - ctx: Context for cancellation and deadlines (not used in current implementation).
//   - id: The unique identifier of the metric to retrieve.
//
// Returns:
//   - A pointer to the found metric, or nil if no metric with the given ID exists.
//   - An error if retrieval fails (currently always nil as no error handling is implemented).
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
