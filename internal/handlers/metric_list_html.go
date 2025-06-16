package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

type MetricHTMLLister interface {
	List(ctx context.Context) ([]types.Metrics, error)
}

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

func handleMetricListHTMLError(w http.ResponseWriter, err error) {
	switch err {
	default:
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}
}
