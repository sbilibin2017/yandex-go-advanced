package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricBodyGetter defines the interface for fetching a metric by its ID
// from a data source or service.
type MetricBodyGetter interface {
	// Get retrieves a metric given its ID.
	// Returns the metric or an error if retrieval fails.
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

// NewMetricGetBodyHandler creates an HTTP handler function that processes
// metric retrieval requests with the metric ID provided in the request body as JSON.
//
// The handler decodes the request body into a MetricID, validates it,
// retrieves the metric from the provided service, and returns the metric as JSON.
//
// Parameters:
//   - val: A validation function to verify the metric ID.
//   - svc: A service implementing MetricBodyGetter to fetch the metric.
//
// Returns:
//   - An http.HandlerFunc that can be registered with an HTTP server.
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

// handleMetricGetBodyError writes appropriate HTTP error responses
// based on the provided error when processing a metric get request.
//
// It distinguishes between invalid metric IDs, not found errors,
// invalid metric types, and internal server errors.
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
