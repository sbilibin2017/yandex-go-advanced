package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewMetricRouter(
	metricUpdatePathHandler http.HandlerFunc,
	metricGetPathHandler http.HandlerFunc,
	metricListHTMLHandler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middlewares...)

	router.Post("/update/{type}/{name}/{value}", metricUpdatePathHandler)
	router.Get("/value/{type}/{name}", metricGetPathHandler)
	router.Get("/", metricListHTMLHandler)

	return router
}
