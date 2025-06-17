package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricPathGetter defines the interface for fetching a metric by its ID
// from a data source or service.
type MetricPathGetter interface {
	// Get retrieves a metric given its ID.
	// Returns the metric or an error if retrieval fails.
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

// NewMetricGetPathHandler creates an HTTP handler function that processes
// metric retrieval requests with the metric type and name provided as URL path parameters.
//
// The handler extracts the "type" and "name" parameters from the URL path,
// validates them using the provided validation function,
// fetches the metric from the service, and writes the metric's string value
// in the response body.
//
// Parameters:
//   - val: A validation function that validates metric type and name strings.
//   - svc: A service implementing MetricPathGetter to fetch the metric.
//
// Returns:
//   - An http.HandlerFunc that can be registered with an HTTP server/router.
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

// handleMetricGetPathError writes appropriate HTTP error responses based on the
// provided error when processing a metric get request from the URL path.
//
// It distinguishes between missing metric names, not found errors,
// invalid metric types, and internal server errors.
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
