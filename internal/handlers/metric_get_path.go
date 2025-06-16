package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricPathGetter interface {
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

func NewMetricGetPathHandler(
	val func(metricType string, metricName string) error,
	svc MetricPathGetter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")

		err := val(metricType, metricName)
		if err != nil {
			handleMetricGetPathError(w, err)
			return
		}

		id := types.NewMetricID(metricType, metricName)

		metric, err := svc.Get(r.Context(), *id)
		if err != nil {
			handleMetricGetPathError(w, err)
			return
		}
		if metric == nil {
			handleMetricGetPathError(w, errors.ErrMetricNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(types.GetMetricsStringValue(metric)))

	}
}

func handleMetricGetPathError(w http.ResponseWriter, err error) {
	switch err {
	case errors.ErrMetricNameMissing, errors.ErrMetricNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.ErrMetricTypeInvalid:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}
}
