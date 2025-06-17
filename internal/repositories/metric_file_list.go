package repositories

import (
	"context"
	"encoding/json"
	"os"
	"sort"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricFileListRepository struct {
	file string
}

func NewMetricFileListRepository(filePath string) *MetricFileListRepository {
	return &MetricFileListRepository{
		file: filePath,
	}
}

func (repo *MetricFileListRepository) List(ctx context.Context) ([]types.Metrics, error) {
	mu.RLock()
	defer mu.RUnlock()

	f, err := os.Open(repo.file)
	if err != nil {
		if os.IsNotExist(err) {
			return []types.Metrics{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var list []types.Metrics
	decoder := json.NewDecoder(f)

	for {
		var m types.Metrics
		err := decoder.Decode(&m)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		list = append(list, m)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})

	return list, nil
}
