package types_test

import (
	"testing"

	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricID(t *testing.T) {
	tests := []struct {
		metricType string
		metricName string
		wantID     string
		wantType   string
	}{
		{"gauge", "temperature", "temperature", "gauge"},
		{"counter", "requests", "requests", "counter"},
	}

	for _, tt := range tests {
		t.Run(tt.metricType+"-"+tt.metricName, func(t *testing.T) {
			id := types.NewMetricID(tt.metricType, tt.metricName)
			assert.Equal(t, tt.wantID, id.ID)
			assert.Equal(t, tt.wantType, id.Type)
		})
	}
}

func TestNewMetric(t *testing.T) {
	tests := []struct {
		metricType  string
		metricName  string
		metricValue string
		wantDelta   *int64
		wantValue   *float64
	}{
		{
			metricType:  types.Gauge,
			metricName:  "temp",
			metricValue: "123.45",
			wantDelta:   nil,
			wantValue:   ptrFloat64(123.45),
		},
		{
			metricType:  types.Counter,
			metricName:  "requests",
			metricValue: "100",
			wantDelta:   ptrInt64(100),
			wantValue:   nil,
		},
		{
			metricType:  types.Gauge,
			metricName:  "temp",
			metricValue: "abc",
			wantDelta:   nil,
			wantValue:   nil,
		},
		{
			metricType:  types.Counter,
			metricName:  "requests",
			metricValue: "abc",
			wantDelta:   nil,
			wantValue:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.metricType+"-"+tt.metricName, func(t *testing.T) {
			m := types.NewMetric(tt.metricType, tt.metricName, tt.metricValue)
			assert.Equal(t, tt.metricName, m.ID)
			assert.Equal(t, tt.metricType, m.Type)

			if tt.wantDelta == nil {
				assert.Nil(t, m.Delta)
			} else {
				assert.NotNil(t, m.Delta)
				assert.Equal(t, *tt.wantDelta, *m.Delta)
			}

			if tt.wantValue == nil {
				assert.Nil(t, m.Value)
			} else {
				assert.NotNil(t, m.Value)
				assert.InDelta(t, *tt.wantValue, *m.Value, 0.00001)
			}
		})
	}
}

func TestNewMetricsHTML(t *testing.T) {
	gv := 12.34
	cv := int64(56)

	tests := []struct {
		name     string
		metrics  []types.Metrics
		expected []string
	}{
		{
			name: "Mixed metrics",
			metrics: []types.Metrics{
				{ID: "gauge1", Type: types.Gauge, Value: &gv},
				{ID: "counter1", Type: types.Counter, Delta: &cv},
				{ID: "unknown", Type: "unknown"},
				{ID: "gauge_nil", Type: types.Gauge},
				{ID: "counter_nil", Type: types.Counter},
			},
			expected: []string{
				"<li>gauge1: 12.34</li>",
				"<li>counter1: 56</li>",
				"<li>unknown: N/A</li>",
				"<li>gauge_nil: N/A</li>",
				"<li>counter_nil: N/A</li>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := types.NewMetricsHTML(tt.metrics)
			assert.Contains(t, html, "<html>")
			for _, exp := range tt.expected {
				assert.Contains(t, html, exp)
			}
			assert.Contains(t, html, "</html>")
		})
	}
}

func TestGetMetricsStringValue(t *testing.T) {
	gv := 78.9
	cv := int64(123)

	tests := []struct {
		name     string
		metric   *types.Metrics
		expected string
	}{
		{"Gauge with value", &types.Metrics{Type: types.Gauge, Value: &gv}, "78.9"},
		{"Counter with delta", &types.Metrics{Type: types.Counter, Delta: &cv}, "123"},
		{"Gauge with nil value", &types.Metrics{Type: types.Gauge}, ""},
		{"Counter with nil delta", &types.Metrics{Type: types.Counter}, ""},
		{"Unknown type", &types.Metrics{Type: "unknown"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.GetMetricsStringValue(tt.metric)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// helpers to create pointers

func ptrInt64(v int64) *int64 {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}
