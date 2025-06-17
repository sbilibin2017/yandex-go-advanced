package validators

import (
	"strconv"

	"github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

// ValidateMetricIDAttributes checks if the provided metricType and metricName
// are valid. Returns an error if the metricName is empty or the metricType
// is not one of the recognized types (Counter or Gauge).
func ValidateMetricIDAttributes(metricType, metricName string) error {
	if metricName == "" {
		return errors.ErrMetricNameMissing
	}

	if metricType != types.Counter && metricType != types.Gauge {
		return errors.ErrMetricTypeInvalid
	}

	return nil
}

// ValidateMetricAttributes validates the metricType, metricName, and metricValue.
// It first validates the metric ID attributes, then checks if the metricValue
// can be properly parsed according to the metricType.
// Returns an error if any validation fails.
func ValidateMetricAttributes(metricType, metricName, metricValue string) error {
	err := ValidateMetricIDAttributes(metricType, metricName)
	if err != nil {
		return err
	}

	switch metricType {
	case types.Counter:
		_, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.ErrMetricValueInvalid
		}
	case types.Gauge:
		_, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors.ErrMetricValueInvalid
		}
	}

	return nil
}

// ValidateMetricID validates a MetricID struct, ensuring its ID and Type are valid.
// Returns an error if ID is empty or Type is not recognized.
func ValidateMetricID(id types.MetricID) error {
	if id.ID == "" {
		return errors.ErrMetricIDInvalid
	}

	if id.Type != types.Counter && id.Type != types.Gauge {
		return errors.ErrMetricTypeInvalid
	}

	return nil
}

// ValidateMetric validates a Metrics struct, ensuring its ID and type are valid
// and that the corresponding value field is set (Delta for Counter, Value for Gauge).
// Returns an error if any validation fails.
func ValidateMetric(metric types.Metrics) error {
	err := ValidateMetricID(types.MetricID{ID: metric.ID, Type: metric.Type})
	if err != nil {
		return err
	}

	switch metric.Type {
	case types.Counter:
		if metric.Delta == nil {
			return errors.ErrMetricDeltaInvalid
		}
	case types.Gauge:
		if metric.Value == nil {
			return errors.ErrMetricValueInvalid
		}
	}

	return nil
}
