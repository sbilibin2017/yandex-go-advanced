package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	internalErrors "github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricUpdaterPath interface {
	Update(ctx context.Context, metrics []*types.Metrics) ([]*types.Metrics, error)
}

func NewMetricUpdatePathHandler(
	val func(metricType string, metricName string, metricValue string) error,
	svc MetricUpdaterPath,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		err := val(metricType, metricName, metricValue)
		if err != nil {
			handleMetricUpdatePathError(w, err)
			return
		}

		metric := types.NewMetric(metricType, metricName, metricValue)

		// Передаём срез указателей, а не значений
		_, err = svc.Update(r.Context(), []*types.Metrics{metric})
		if err != nil {
			handleMetricUpdatePathError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleMetricUpdatePathError(w http.ResponseWriter, err error) {
	switch err {
	case internalErrors.ErrMetricNameMissing:
		http.Error(w, err.Error(), http.StatusNotFound)
	case internalErrors.ErrMetricTypeInvalid,
		internalErrors.ErrMetricValueInvalid:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, internalErrors.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}
}
