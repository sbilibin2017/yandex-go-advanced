package services

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricUpdateSaver interface {
	Save(ctx context.Context, metrics types.Metrics) error
}

type MetricUpdateGetter interface {
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

type MetricUpdateService struct {
	saver  MetricUpdateSaver
	getter MetricUpdateGetter
}

func NewMetricUpdateService(
	saver MetricUpdateSaver,
	getter MetricUpdateGetter,
) *MetricUpdateService {
	return &MetricUpdateService{saver: saver, getter: getter}
}

func (svc *MetricUpdateService) Update(
	ctx context.Context,
	metrics []*types.Metrics,
) ([]*types.Metrics, error) {
	for idx, m := range metrics {
		switch m.Type {
		case types.Counter:
			if err := updateCounterMetric(ctx, svc.getter, m); err != nil {
				return nil, err
			}
		}

		err := svc.saver.Save(ctx, *m)
		if err != nil {
			return nil, err
		}
		metrics[idx] = m
	}

	return metrics, nil
}

func updateCounterMetric(
	ctx context.Context,
	getter MetricUpdateGetter,
	metric *types.Metrics,
) error {
	existing, err := getter.Get(ctx, types.MetricID{ID: metric.ID, Type: metric.Type})
	if err != nil {
		return err
	}

	if existing != nil && existing.Delta != nil && metric.Delta != nil {
		*metric.Delta += *existing.Delta
	}

	return nil
}
