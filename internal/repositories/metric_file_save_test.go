package repositories

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

func TestMetricFileSaveRepository_Save(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "metrics.jsonl")

	repo := NewMetricFileSaveRepository(filePath)
	ctx := context.Background()

	// Помощники для значений
	fval := func(v float64) *float64 { return &v }
	ival := func(v int64) *int64 { return &v }

	metric1 := types.Metrics{
		ID:    "metric1",
		Type:  "gauge",
		Value: fval(3.14),
	}
	metric2 := types.Metrics{
		ID:    "metric2",
		Type:  "counter",
		Delta: ival(10),
	}
	metric1Updated := types.Metrics{
		ID:    "metric1",
		Type:  "gauge",
		Value: fval(2.71),
	}

	// 1. Сохраняем metric1
	err := repo.Save(ctx, metric1)
	require.NoError(t, err)

	// Проверяем файл: должна быть одна запись, равная metric1
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)

	var readMetrics []types.Metrics
	for _, line := range splitLines(data) {
		var m types.Metrics
		err := json.Unmarshal(line, &m)
		require.NoError(t, err)
		readMetrics = append(readMetrics, m)
	}
	assert.Len(t, readMetrics, 1)
	assert.Equal(t, metric1, readMetrics[0])

	// 2. Добавляем metric2
	err = repo.Save(ctx, metric2)
	require.NoError(t, err)

	data, err = os.ReadFile(filePath)
	require.NoError(t, err)

	readMetrics = nil
	for _, line := range splitLines(data) {
		var m types.Metrics
		err := json.Unmarshal(line, &m)
		require.NoError(t, err)
		readMetrics = append(readMetrics, m)
	}
	assert.Len(t, readMetrics, 2)

	// 3. Обновляем metric1
	err = repo.Save(ctx, metric1Updated)
	require.NoError(t, err)

	data, err = os.ReadFile(filePath)
	require.NoError(t, err)

	readMetrics = nil
	for _, line := range splitLines(data) {
		var m types.Metrics
		err := json.Unmarshal(line, &m)
		require.NoError(t, err)
		readMetrics = append(readMetrics, m)
	}

	assert.Len(t, readMetrics, 2)
	foundMetric1 := false
	for _, m := range readMetrics {
		if m.ID == metric1Updated.ID && m.Type == metric1Updated.Type {
			assert.NotNil(t, m.Value)
			assert.Equal(t, *metric1Updated.Value, *m.Value)
			foundMetric1 = true
		}
	}
	assert.True(t, foundMetric1)
}

// splitLines разбивает []byte на слайс строк без переноса строки
func splitLines(data []byte) [][]byte {
	lines := [][]byte{}
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	// Последняя строка, если нет перевода строки в конце
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
