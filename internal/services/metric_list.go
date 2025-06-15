package services

import (
	"context"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricListLister interface {
	List(ctx context.Context) ([]types.Metrics, error)
}

type MetricListService struct {
	lister MetricListLister
}

func NewMetricListService(
	lister MetricListLister,
) *MetricListService {
	return &MetricListService{lister: lister}
}

func (svc *MetricListService) List(
	ctx context.Context,
) ([]types.Metrics, error) {
	metrics, err := svc.lister.List(ctx)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
