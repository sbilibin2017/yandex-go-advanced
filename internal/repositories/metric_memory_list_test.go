package repositories

import (
	"context"
	"sort"
	"testing"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricMemoryListRepository_List(t *testing.T) {
	// Очищаем глобальное хранилище
	mu.Lock()
	metrics = make(map[types.MetricID]types.Metrics)
	mu.Unlock()

	repo := NewMetricMemoryListRepository()
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

	mu.Lock()
	metrics[types.MetricID{ID: metric1.ID, MType: metric1.MType}] = metric1
	metrics[types.MetricID{ID: metric2.ID, MType: metric2.MType}] = metric2
	metrics[types.MetricID{ID: metric3.ID, MType: metric3.MType}] = metric3
	mu.Unlock()

	tests := []struct {
		name    string
		want    []types.Metrics
		wantErr bool
	}{
		{
			name: "list all metrics sorted by ID",
			want: []types.Metrics{metric2, metric1, metric3}, // metricA, metricB, metricC
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Проверка сортировки по ID
			assert.Equal(t, tt.want, got)

			// Альтернатива: если порядок важен, можно отсортировать вручную и сравнить
			sortedGot := append([]types.Metrics(nil), got...)
			sort.Slice(sortedGot, func(i, j int) bool {
				return sortedGot[i].ID < sortedGot[j].ID
			})
			assert.Equal(t, tt.want, sortedGot)
		})
	}
}
