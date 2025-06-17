package repositories

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricFileSaveRepository struct {
	file string
}

func NewMetricFileSaveRepository(filePath string) *MetricFileSaveRepository {
	return &MetricFileSaveRepository{
		file: filePath,
	}
}

func (repo *MetricFileSaveRepository) Save(ctx context.Context, m types.Metrics) error {
	mu.Lock()
	defer mu.Unlock()

	origFile, err := os.Open(repo.file)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	defer func() {
		if origFile != nil {
			origFile.Close()
		}
	}()

	tempFilePath := repo.file + ".tmp"
	tempFile, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	found := false

	if origFile != nil {
		scanner := bufio.NewScanner(origFile)
		for scanner.Scan() {
			line := scanner.Bytes()

			var existingMetric types.Metrics
			if err := json.Unmarshal(line, &existingMetric); err != nil {
				if _, err := tempFile.Write(line); err != nil {
					return err
				}
				if _, err := tempFile.Write([]byte("\n")); err != nil {
					return err
				}
				continue
			}

			if existingMetric.ID == m.ID && existingMetric.Type == m.Type {
				data, err := json.Marshal(m)
				if err != nil {
					return err
				}
				if _, err := tempFile.Write(data); err != nil {
					return err
				}
				if _, err := tempFile.Write([]byte("\n")); err != nil {
					return err
				}
				found = true
			} else {
				if _, err := tempFile.Write(line); err != nil {
					return err
				}
				if _, err := tempFile.Write([]byte("\n")); err != nil {
					return err
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	if !found {
		data, err := json.Marshal(m)
		if err != nil {
			return err
		}
		if _, err := tempFile.Write(data); err != nil {
			return err
		}
		if _, err := tempFile.Write([]byte("\n")); err != nil {
			return err
		}
	}

	if err := tempFile.Close(); err != nil {
		return err
	}
	if origFile != nil {
		origFile.Close()
	}

	if err := os.Rename(tempFilePath, repo.file); err != nil {
		return err
	}

	return nil
}
