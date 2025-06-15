package services

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricGetGetter interface {
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

type MetricGetService struct {
	getter MetricGetGetter
}

func NewMetricGetService(
	getter MetricGetGetter,
) *MetricGetService {
	return &MetricGetService{getter: getter}
}

func (svc *MetricGetService) Get(
	ctx context.Context,
	id types.MetricID,
) (*types.Metrics, error) {
	metric, err := svc.getter.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return metric, nil
}
