package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricUpdaterBody interface {
	Update(ctx context.Context, metrics []*types.Metrics) ([]*types.Metrics, error)
}

func NewMetricUpdateBodyHandler(
	val func(metric types.Metrics) error,
	svc MetricUpdaterBody,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var metric types.Metrics

		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			log.Printf("failed to decode JSON body: %v", err)
			http.Error(w, "invalid JSON format", http.StatusBadRequest)
			return
		}

		err = val(metric)
		if err != nil {
			handleMetricUpdateBodyError(w, err)
			return
		}

		metrics, err := svc.Update(r.Context(), []*types.Metrics{&metric})
		if err != nil {
			handleMetricUpdateBodyError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(*metrics[0])

	}
}

func handleMetricUpdateBodyError(w http.ResponseWriter, err error) {
	switch err {
	case errors.ErrMetricIDInvalid:
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.ErrMetricTypeInvalid,
		errors.ErrMetricDeltaInvalid,
		errors.ErrMetricValueInvalid:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}
}
