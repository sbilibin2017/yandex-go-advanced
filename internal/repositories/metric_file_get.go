package repositories

import (
	"context"
	"encoding/json"
	"os"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricFileGetRepository struct {
	file string
}

func NewMetricFileGetRepository(filePath string) *MetricFileGetRepository {
	return &MetricFileGetRepository{
		file: filePath,
	}
}

func (repo *MetricFileGetRepository) Get(
	ctx context.Context,
	id types.MetricID,
) (*types.Metrics, error) {
	mu.RLock()
	defer mu.RUnlock()

	f, err := os.Open(repo.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

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

		if m.ID == id.ID && m.Type == id.Type {
			return &m, nil
		}
	}

	return nil, nil
}
