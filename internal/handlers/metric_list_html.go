package handlers

import (
	"context"
	"html"
	"net/http"
	"strconv"

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
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(newMetricsHTML(metrics)))
	}
}

func newMetricsHTML(metrics []types.Metrics) string {
	htmlStr := "<html><head><title>Metrics List</title></head><body>"
	htmlStr += "<h1>Metrics</h1>"
	htmlStr += "<ul>"

	for _, metric := range metrics {
		name := html.EscapeString(metric.ID)
		var value string

		switch metric.MType {
		case types.Gauge:
			value = metricListHTMLFormatGauge(metric.Value)
		case types.Counter:
			value = metricListHTMLFormatCounter(metric.Delta)
		default:
			value = ""
		}

		value = html.EscapeString(value)
		htmlStr += "<li>" + name + ": " + value + "</li>"
	}

	htmlStr += "</ul></body></html>"
	return htmlStr
}

func metricListHTMLFormatGauge(val *float64) string {
	if val == nil {
		return "0.0"
	}
	return strconv.FormatFloat(*val, 'f', -1, 64)
}

func metricListHTMLFormatCounter(delta *int64) string {
	if delta == nil {
		return "0"
	}
	return strconv.FormatInt(*delta, 10)
}
