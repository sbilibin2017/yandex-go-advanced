package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

func TestMetricMemorySaveRepository_Save(t *testing.T) {
	// Очищаем глобальное хранилище перед запуском теста
	mu.Lock()
	metrics = make(map[types.MetricID]types.Metrics)
	mu.Unlock()

	repo := NewMetricMemorySaveRepository()

	ptrInt64 := func(i int64) *int64 { return &i }
	ptrFloat64 := func(f float64) *float64 { return &f }

	ctx := context.Background()

	tests := []struct {
		name  string
		input types.Metrics
	}{
		{
			name: "save counter metric",
			input: types.Metrics{
				ID:    "metric1",
				MType: types.Counter,
				Delta: ptrInt64(10),
			},
		},
		{
			name: "save gauge metric",
			input: types.Metrics{
				ID:    "metric2",
				MType: types.Gauge,
				Value: ptrFloat64(3.14),
			},
		},
		{
			name: "save metric without optional values",
			input: types.Metrics{
				ID:    "metric3",
				MType: types.Gauge,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(ctx, tt.input)
			assert.NoError(t, err)

			key := types.MetricID{ID: tt.input.ID, MType: tt.input.MType}

			mu.RLock()
			savedMetric, ok := metrics[key]
			mu.RUnlock()

			assert.True(t, ok, "metric should be saved in global metrics")
			assert.Equal(t, tt.input, savedMetric)
		})
	}
}
