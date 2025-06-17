package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// MetricHTMLLister defines the interface for listing metrics as a slice.
// Implementations should provide a method to retrieve all metrics.
type MetricHTMLLister interface {
	// List retrieves all available metrics.
	// Returns a slice of Metrics or an error if retrieval fails.
	List(ctx context.Context) ([]types.Metrics, error)
}

// NewMetricListHTMLHandler returns an HTTP handler function that
// serves an HTML page listing all metrics.
//
// It fetches the metrics from the provided MetricHTMLLister service,
// sets the appropriate Content-Type header, and writes the HTML response.
//
// Parameters:
//   - svc: a service implementing MetricHTMLLister to fetch the metrics.
//
// Returns:
//   - http.HandlerFunc that can be registered to serve metric listings.
func NewMetricListHTMLHandler(
	svc MetricHTMLLister,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := svc.List(r.Context())
		if err != nil {
			handleMetricListHTMLError(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(types.NewMetricsHTML(metrics)))
	}
}

// handleMetricListHTMLError handles errors that occur during metric listing
// by sending an HTTP 500 Internal Server Error response with a generic message.
//
// This function can be extended to handle more specific errors as needed.
func handleMetricListHTMLError(w http.ResponseWriter, err error) {
	switch err {
	default:
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}
}
