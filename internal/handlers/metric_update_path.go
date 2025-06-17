package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	internalErrors "github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricUpdaterPath defines an interface for updating metrics from URL path parameters.
// Implementations should handle the update logic and return the updated metrics or an error.
type MetricUpdaterPath interface {
	// Update processes and updates the given slice of metric pointers.
	// Returns the updated slice of metrics or an error if the update fails.
	Update(ctx context.Context, metrics []*types.Metrics) ([]*types.Metrics, error)
}

// NewMetricUpdatePathHandler returns an HTTP handler function that updates a metric
// based on the type, name, and value extracted from URL path parameters.
//
// It validates the metric attributes using the provided validation function,
// constructs a metric instance, and calls the update service to apply the update.
//
// Parameters:
//   - val: validation function that checks the metric type, name, and value.
//   - svc: service implementing MetricUpdaterPath to perform the update.
//
// Returns:
//   - http.HandlerFunc to be used as an HTTP handler.
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

		// Pass a slice of metric pointers to the update service.
		_, err = svc.Update(r.Context(), []*types.Metrics{metric})
		if err != nil {
			handleMetricUpdatePathError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// handleMetricUpdatePathError writes an appropriate HTTP error response
// depending on the error type encountered during metric update via path parameters.
//
// Known errors map to specific HTTP status codes, and unknown errors
// result in a 500 Internal Server Error response.
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
