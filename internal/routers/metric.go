package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NewMetricRouter creates and returns a new HTTP router configured with routes
// for metric updates, retrievals, and listing, along with optional middleware.
//
// Parameters:
//   - metricUpdatePathHandler: Handler for metric updates via URL path parameters.
//   - metricUpdateBodyHandler: Handler for metric updates via JSON body.
//   - metricGetPathHandler: Handler for metric retrieval via URL path parameters.
//   - metricGetBodyHandler: Handler for metric retrieval via JSON body.
//   - metricListHTMLHandler: Handler for listing all metrics as HTML.
//   - middlewares: Optional variadic middleware functions applied to all routes.
//
// Returns:
//   - An http.Handler that routes requests to the appropriate metric handlers.
func NewMetricRouter(
	metricUpdatePathHandler http.HandlerFunc,
	metricUpdateBodyHandler http.HandlerFunc,
	metricGetPathHandler http.HandlerFunc,
	metricGetBodyHandler http.HandlerFunc,
	metricListHTMLHandler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middlewares...)

	router.Post("/update/{type}/{name}/{value}", metricUpdatePathHandler)
	router.Post("/update/", metricUpdateBodyHandler)

	router.Get("/value/{type}/{name}", metricGetPathHandler)
	router.Post("/value/", metricGetBodyHandler)

	router.Get("/", metricListHTMLHandler)

	return router
}
