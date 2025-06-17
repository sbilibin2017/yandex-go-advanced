package repositories_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sbilibin2017/yandex-go-advanced/internal/repositories"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

func TestMetricFileListRepository_List(t *testing.T) {
	int64ptr := func(i int64) *int64 { return &i }
	float64ptr := func(f float64) *float64 { return &f }

	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Запишем несколько метрик построчно (каждая — отдельный JSON)
	metrics := []types.Metrics{
		{ID: "metric1", Type: "counter", Delta: int64ptr(100)},
		{ID: "metric2", Type: "gauge", Value: float64ptr(12.34)},
		{ID: "metric0", Type: "counter", Delta: int64ptr(50)},
	}

	for _, m := range metrics {
		data, err := json.Marshal(m)
		require.NoError(t, err)

		_, err = tmpFile.Write(data)
		require.NoError(t, err)

		_, err = tmpFile.Write([]byte("\n"))
		require.NoError(t, err)
	}

	require.NoError(t, tmpFile.Close())

	repo := repositories.NewMetricFileListRepository(tmpFile.Name())

	list, err := repo.List(context.Background())
	require.NoError(t, err)

	// Проверяем, что возвращаемые метрики отсортированы по ID
	assert.Len(t, list, len(metrics))
	assert.Equal(t, "metric0", list[0].ID)
	assert.Equal(t, "metric1", list[1].ID)
	assert.Equal(t, "metric2", list[2].ID)
}
