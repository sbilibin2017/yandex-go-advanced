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
