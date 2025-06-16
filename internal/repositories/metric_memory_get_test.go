package repositories

import (
	"context"
	"testing"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricMemoryGetRepository_Get(t *testing.T) {
	// Очищаем глобальное хранилище перед началом теста
	mu.Lock()
	metrics = make(map[types.MetricID]types.Metrics)
	mu.Unlock()

	repo := NewMetricMemoryGetRepository()
	ctx := context.Background()

	ptrInt64 := func(i int64) *int64 {
		return &i
	}

	// Подготовка данных
	existingMetric := types.Metrics{
		ID:    "metric1",
		Type:  types.Counter,
		Delta: ptrInt64(42),
	}
	key := types.MetricID{ID: existingMetric.ID, Type: existingMetric.Type}

	mu.Lock()
	metrics[key] = existingMetric
	mu.Unlock()

	tests := []struct {
		name    string
		id      types.MetricID
		want    *types.Metrics
		wantErr bool
	}{
		{
			name: "found existing metric",
			id:   key,
			want: &existingMetric,
		},
		{
			name:    "metric not found",
			id:      types.MetricID{ID: "not_exist", Type: types.Gauge},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
