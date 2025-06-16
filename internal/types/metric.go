package types

import (
	"html"
	"strconv"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type MetricID struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func NewMetricID(
	metricType string,
	metricName string,
) *MetricID {
	return &MetricID{
		ID:   metricName,
		Type: metricType,
	}
}

type Metrics struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func NewMetric(
	metricType string,
	metricName string,
	metricValue string,
) *Metrics {
	m := &Metrics{
		ID:   metricName,
		Type: metricType,
	}

	switch metricType {
	case Gauge:
		if val, err := strconv.ParseFloat(metricValue, 64); err == nil {
			m.Value = &val
		}
	case Counter:
		if val, err := strconv.ParseInt(metricValue, 10, 64); err == nil {
			m.Delta = &val
		}
	}

	return m
}

func NewMetricsHTML(metrics []Metrics) string {
	htmlStr := "<html><head><title>Metrics List</title></head><body>"
	htmlStr += "<h1>Metrics</h1>"
	htmlStr += "<ul>"

	for _, metric := range metrics {
		name := html.EscapeString(metric.ID)
		var value string

		switch metric.Type {
		case Gauge:
			if metric.Value == nil {
				value = "N/A"
			} else {
				value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
			}
		case Counter:
			if metric.Delta == nil {
				value = "N/A"
			} else {
				value = strconv.FormatInt(*metric.Delta, 10)
			}
		default:
			value = "N/A"
		}

		value = html.EscapeString(value)
		htmlStr += "<li>" + name + ": " + value + "</li>"
	}

	htmlStr += "</ul></body></html>"
	return htmlStr
}

func GetMetricsStringValue(metric *Metrics) string {
	switch metric.Type {
	case Gauge:
		if metric.Value != nil {
			return strconv.FormatFloat(*metric.Value, 'f', -1, 64)
		}
	case Counter:
		if metric.Delta != nil {
			return strconv.FormatInt(*metric.Delta, 10)
		}
	}
	return ""
}
