package types

import (
	"html"
	"strconv"
)

const (
	// Counter represents a metric type for integer counters.
	Counter = "counter"
	// Gauge represents a metric type for floating-point gauges.
	Gauge = "gauge"
)

// MetricID uniquely identifies a metric by its type and name.
type MetricID struct {
	ID   string `json:"id"`   // Metric name or identifier
	Type string `json:"type"` // Metric type, e.g. "counter" or "gauge"
}

// NewMetricID creates a new MetricID given a metric type and name.
func NewMetricID(
	metricType string,
	metricName string,
) *MetricID {
	return &MetricID{
		ID:   metricName,
		Type: metricType,
	}
}

// Metrics represents a metric with its ID, type, and value(s).
// For counter metrics, Delta holds an integer count.
// For gauge metrics, Value holds a floating-point measurement.
// Hash can be used for data integrity or verification.
type Metrics struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

// NewMetric constructs a new Metrics instance from the provided type, name, and value string.
// It parses the value string into the appropriate type based on metricType.
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

// NewMetricsHTML generates an HTML page listing the provided metrics.
// Each metric is escaped properly to prevent HTML injection.
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

// GetMetricsStringValue returns the string representation of a metric's value.
// Returns an empty string if the value is not set or the metric type is unknown.
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
