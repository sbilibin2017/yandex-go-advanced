package repositories

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricMemorySaveRepository provides an in-memory repository for saving metrics.
type MetricMemorySaveRepository struct{}

// NewMetricMemorySaveRepository creates and returns a new MetricMemorySaveRepository instance.
func NewMetricMemorySaveRepository() *MetricMemorySaveRepository {
	return &MetricMemorySaveRepository{}
}

// Save stores the given metric in memory, keyed by its MetricID.
//
// Parameters:
//   - ctx: Context for cancellation and deadlines (not used in current implementation).
//   - m: The Metrics value to save.
//
// Returns:
//   - An error if saving fails (currently always nil as no error handling is implemented).
func (repo *MetricMemorySaveRepository) Save(
	ctx context.Context,
	m types.Metrics,
) error {
	mu.Lock()
	defer mu.Unlock()

	metrics[types.MetricID{ID: m.ID, Type: m.Type}] = m
	return nil
}
