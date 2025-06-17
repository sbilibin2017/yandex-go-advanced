package repositories_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sbilibin2017/yandex-go-advanced/internal/repositories"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

func TestMetricFileGetRepository_Get(t *testing.T) {

	ptrInt64 := func(v int64) *int64 {
		return &v
	}

	ptrFloat64 := func(v float64) *float64 {
		return &v
	}

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "metrics.json")

	metrics := []types.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: ptrInt64(42),
		},
		{
			ID:    "metric2",
			Type:  "gauge",
			Value: ptrFloat64(3.14),
		},
	}

	// Записать метрики в файл по одной JSON-строке на строку
	f, err := os.Create(filePath)
	require.NoError(t, err)
	for _, m := range metrics {
		err = json.NewEncoder(f).Encode(m)
		require.NoError(t, err)
	}
	f.Close()

	repo := repositories.NewMetricFileGetRepository(filePath)
	ctx := context.Background()

	t.Run("get existing counter metric", func(t *testing.T) {
		m, err := repo.Get(ctx, types.MetricID{ID: "metric1", Type: "counter"})
		require.NoError(t, err)
		require.NotNil(t, m)
		assert.Equal(t, "metric1", m.ID)
		assert.Equal(t, "counter", m.Type)
		assert.NotNil(t, m.Delta)
		assert.Equal(t, int64(42), *m.Delta)
		assert.Nil(t, m.Value)
	})

	t.Run("get existing gauge metric", func(t *testing.T) {
		m, err := repo.Get(ctx, types.MetricID{ID: "metric2", Type: "gauge"})
		require.NoError(t, err)
		require.NotNil(t, m)
		assert.Equal(t, "metric2", m.ID)
		assert.Equal(t, "gauge", m.Type)
		assert.NotNil(t, m.Value)
		assert.Equal(t, 3.14, *m.Value)
		assert.Nil(t, m.Delta)
	})

	t.Run("get non-existing metric", func(t *testing.T) {
		m, err := repo.Get(ctx, types.MetricID{ID: "notfound", Type: "counter"})
		require.NoError(t, err)
		assert.Nil(t, m)
	})

	t.Run("file does not exist", func(t *testing.T) {
		repo := repositories.NewMetricFileGetRepository(filepath.Join(tmpDir, "no_file.json"))
		m, err := repo.Get(ctx, types.MetricID{ID: "metric1", Type: "counter"})
		require.NoError(t, err)
		assert.Nil(t, m)
	})
}
