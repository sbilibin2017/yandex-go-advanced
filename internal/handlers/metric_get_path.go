package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricPathGetter interface {
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

func NewMetricGetPathHandler(svc MetricPathGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")

		id, err := validateMetricGetPath(metricType, metricName)
		if err != nil {
			handleMetricGetPathError(w, err)
			return
		}

		metric, err := svc.Get(r.Context(), id)
		if err != nil {
			handleMetricGetPathError(w, err)
			return
		}
		if metric == nil {
			handleMetricGetPathError(w, errMetricGetNotFound)
			return
		}

		switch metric.MType {
		case types.Gauge:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(metricGetPathFormatGauge(metric.Value)))
		case types.Counter:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(metricGetPathFormatCounter(metric.Delta)))
		default:
			handleMetricGetPathError(w, errMetricGetTypeInvalid)
		}
	}
}

var (
	errMetricGetNameMissing = errors.New("metric name is required")
	errMetricGetTypeInvalid = errors.New("invalid metric type")
	errMetricGetNotFound    = errors.New("metric not found")
)

func validateMetricGetPath(metricType, metricName string) (types.MetricID, error) {
	if metricName == "" {
		return types.MetricID{}, errMetricGetNameMissing
	}

	var metric types.MetricID
	metric.ID = metricName

	switch metricType {
	case types.Counter:
		metric.MType = types.Counter
	case types.Gauge:
		metric.MType = types.Gauge
	default:
		return types.MetricID{}, errMetricGetTypeInvalid
	}

	return metric, nil
}

func handleMetricGetPathError(w http.ResponseWriter, err error) {
	switch err {
	case errMetricGetNameMissing:
		http.Error(w, err.Error(), http.StatusNotFound)
	case errMetricGetTypeInvalid:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errMetricGetNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func metricGetPathFormatGauge(val *float64) string {
	if val == nil {
		return "0.0"
	}
	return strconv.FormatFloat(*val, 'f', -1, 64)
}

func metricGetPathFormatCounter(delta *int64) string {
	if delta == nil {
		return "0"
	}
	return strconv.FormatInt(*delta, 10)
}
