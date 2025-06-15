package repositories_test

import (
	"context"
	"testing"

	"github.com/sbilibin2017/yandex-go-advanced/internal/repositories"
	"github.com/sbilibin2017/yandex-go-advanced/internal/storages"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricMemoryListRepository_List(t *testing.T) {
	storage := storages.NewMemoryStorage[types.MetricID, types.Metrics]()
	repo := repositories.NewMetricMemoryListRepository(storage)
	ctx := context.Background()

	ptrFloat64 := func(f float64) *float64 {
		return &f
	}

	ptrInt64 := func(i int64) *int64 {
		return &i
	}
	metric1 := types.Metrics{ID: "metricB", MType: types.Gauge, Value: ptrFloat64(3.14)}
	metric2 := types.Metrics{ID: "metricA", MType: types.Counter, Delta: ptrInt64(42)}
	metric3 := types.Metrics{ID: "metricC", MType: types.Gauge, Value: ptrFloat64(2.71)}

	// Наполняем хранилище
	storage.Mu.Lock()
	storage.Store[types.MetricID{ID: metric1.ID, MType: metric1.MType}] = metric1
	storage.Store[types.MetricID{ID: metric2.ID, MType: metric2.MType}] = metric2
	storage.Store[types.MetricID{ID: metric3.ID, MType: metric3.MType}] = metric3
	storage.Mu.Unlock()

	tests := []struct {
		name    string
		want    []types.Metrics
		wantErr bool
	}{
		{
			name: "list all metrics sorted by ID",
			want: []types.Metrics{metric2, metric1, metric3}, // A, B, C
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("List() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
