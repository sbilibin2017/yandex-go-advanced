package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricUpdater interface {
	Update(ctx context.Context, metrics []types.Metrics) error
}

func NewMetricUpdatePathHandler(svc MetricUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		metric, err := validateMetricUpdatePath(metricType, metricName, metricValue)
		if err != nil {
			handleMetricUpdatePathError(w, err)
			return
		}

		err = svc.Update(r.Context(), []types.Metrics{metric})
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

var (
	errMetricNameMissing  = errors.New("metric name is required")
	errMetricTypeInvalid  = errors.New("invalid metric type")
	errMetricValueInvalid = errors.New("invalid metric value")
)

func validateMetricUpdatePath(metricType, metricName, metricValue string) (types.Metrics, error) {
	if metricName == "" {
		return types.Metrics{}, errMetricNameMissing
	}

	var metric types.Metrics
	metric.ID = metricName

	switch metricType {
	case string(types.Counter):
		metric.MType = types.Counter
		delta, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return types.Metrics{}, errMetricValueInvalid
		}
		metric.Delta = &delta

	case string(types.Gauge):
		metric.MType = types.Gauge
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return types.Metrics{}, errMetricValueInvalid
		}
		metric.Value = &value

	default:
		return types.Metrics{}, errMetricTypeInvalid
	}

	return metric, nil
}

func handleMetricUpdatePathError(w http.ResponseWriter, err error) {
	switch err {
	case errMetricNameMissing:
		http.Error(w, err.Error(), http.StatusNotFound)
	case errMetricTypeInvalid, errMetricValueInvalid:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
