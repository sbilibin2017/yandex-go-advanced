package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricUpdaterBody defines an interface for updating metrics from a slice of Metrics pointers.
// Implementations should handle the update logic and return the updated metrics or an error.
type MetricUpdaterBody interface {
	// Update processes and updates the given slice of metric pointers.
	// Returns the updated slice of metrics or an error if the update fails.
	Update(ctx context.Context, metrics []*types.Metrics) ([]*types.Metrics, error)
}

// NewMetricUpdateBodyHandler returns an HTTP handler function that processes
// metric updates sent in the request body as JSON.
//
// It validates the incoming metric using the provided validation function,
// calls the update service to apply the metric update,
// and responds with the updated metric as JSON.
//
// Parameters:
//   - val: a validation function that checks the integrity of the metric.
//   - svc: a service implementing MetricUpdaterBody to perform the update.
//
// Returns:
//   - http.HandlerFunc to be used as an HTTP handler.
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

// handleMetricUpdateBodyError writes an appropriate HTTP error response
// depending on the error type encountered during metric update.
//
// Known errors map to specific HTTP status codes, and unknown errors
// result in a 500 Internal Server Error response.
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
