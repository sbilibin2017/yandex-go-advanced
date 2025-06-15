package repositories_test

import (
	"context"
	"testing"

	"github.com/sbilibin2017/yandex-go-advanced/internal/repositories"
	"github.com/sbilibin2017/yandex-go-advanced/internal/storages"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricMemoryGetRepository_Get(t *testing.T) {
	storage := storages.NewMemoryStorage[types.MetricID, types.Metrics]()
	repo := repositories.NewMetricMemoryGetRepository(storage)
	ctx := context.Background()

	ptrInt64 := func(i int64) *int64 {
		return &i
	}
	// Подготовка данных
	existingMetric := types.Metrics{
		ID:    "metric1",
		MType: types.Counter,
		Delta: ptrInt64(42),
	}
	key := types.MetricID{ID: existingMetric.ID, MType: existingMetric.MType}
	storage.Mu.Lock()
	storage.Store[key] = existingMetric
	storage.Mu.Unlock()

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
			id:      types.MetricID{ID: "not_exist", MType: types.Gauge},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
