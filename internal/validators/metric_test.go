package validators_test

import (
	"testing"

	"github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/sbilibin2017/yandex-go-advanced/internal/validators"
	"github.com/stretchr/testify/assert"
)

func TestValidateMetricIDAttributes(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		metricName string
		expected   error
	}{
		{"missing name", string(types.Gauge), "", errors.ErrMetricNameMissing},
		{"invalid type", "unknown", "test_metric", errors.ErrMetricTypeInvalid},
		{"valid gauge", string(types.Gauge), "temperature", nil},
		{"valid counter", string(types.Counter), "requests", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validators.ValidateMetricIDAttributes(tt.metricType, tt.metricName)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateMetricAttributes(t *testing.T) {
	tests := []struct {
		name        string
		metricType  string
		metricName  string
		metricValue string
		expected    error
	}{
		{"valid gauge", string(types.Gauge), "load", "1.23", nil},
		{"valid counter", string(types.Counter), "hits", "100", nil},
		{"invalid counter value", string(types.Counter), "hits", "abc", errors.ErrMetricValueInvalid},
		{"invalid gauge value", string(types.Gauge), "load", "NaN%", errors.ErrMetricValueInvalid},
		{"missing name", string(types.Gauge), "", "1.5", errors.ErrMetricNameMissing},
		{"invalid type", "unknown", "metric", "42", errors.ErrMetricTypeInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validators.ValidateMetricAttributes(tt.metricType, tt.metricName, tt.metricValue)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestValidateMetricID(t *testing.T) {
	tests := []struct {
		name      string
		metricID  types.MetricID
		wantError error
	}{
		{
			name:      "valid counter",
			metricID:  types.MetricID{ID: "metric1", Type: types.Counter},
			wantError: nil,
		},
		{
			name:      "valid gauge",
			metricID:  types.MetricID{ID: "metric2", Type: types.Gauge},
			wantError: nil,
		},
		{
			name:      "empty ID",
			metricID:  types.MetricID{ID: "", Type: types.Counter},
			wantError: errors.ErrMetricIDInvalid,
		},
		{
			name:      "invalid type",
			metricID:  types.MetricID{ID: "metric3", Type: "invalid"},
			wantError: errors.ErrMetricTypeInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validators.ValidateMetricID(tt.metricID)
			assert.ErrorIs(t, err, tt.wantError)
		})
	}
}

func TestValidateMetric(t *testing.T) {
	delta := int64(10)
	value := float64(3.14)

	tests := []struct {
		name      string
		metric    types.Metrics
		wantError error
	}{
		{
			name: "valid counter metric",
			metric: types.Metrics{
				ID:    "metric1",
				Type:  types.Counter,
				Delta: &delta,
			},
			wantError: nil,
		},
		{
			name: "valid gauge metric",
			metric: types.Metrics{
				ID:    "metric2",
				Type:  types.Gauge,
				Value: &value,
			},
			wantError: nil,
		},
		{
			name: "counter metric with nil delta",
			metric: types.Metrics{
				ID:   "metric3",
				Type: types.Counter,
			},
			wantError: errors.ErrMetricDeltaInvalid,
		},
		{
			name: "gauge metric with nil value",
			metric: types.Metrics{
				ID:   "metric4",
				Type: types.Gauge,
			},
			wantError: errors.ErrMetricValueInvalid,
		},
		{
			name: "invalid metric ID",
			metric: types.Metrics{
				ID:   "",
				Type: types.Counter,
			},
			wantError: errors.ErrMetricIDInvalid,
		},
		{
			name: "invalid metric type",
			metric: types.Metrics{
				ID:   "metric5",
				Type: "invalid",
			},
			wantError: errors.ErrMetricTypeInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validators.ValidateMetric(tt.metric)
			assert.ErrorIs(t, err, tt.wantError)
		})
	}
}
