package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/yandex-go-advanced/internal/errors"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricBodyGetter interface {
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

func NewMetricGetBodyHandler(
	val func(metric types.MetricID) error,
	svc MetricBodyGetter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metricID types.MetricID

		err := json.NewDecoder(r.Body).Decode(&metricID)
		if err != nil {
			http.Error(w, "invalid JSON format", http.StatusBadRequest)
			return
		}

		err = val(metricID)
		if err != nil {
			handleMetricGetBodyError(w, err)
			return
		}

		metric, err := svc.Get(r.Context(), metricID)
		if err != nil {
			handleMetricGetBodyError(w, err)
			return
		}
		if metric == nil {
			handleMetricGetBodyError(w, errors.ErrMetricNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(metric)
	}
}

func handleMetricGetBodyError(w http.ResponseWriter, err error) {
	switch err {
	case errors.ErrMetricIDInvalid, errors.ErrMetricNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.ErrMetricTypeInvalid:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}
}
