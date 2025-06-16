package validators

import (
	"strconv"

	"github.com/sbilibin2017/yandex-go-advanced/internal/errors"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
)

func ValidateMetricIDAttributes(metricType, metricName string) error {
	if metricName == "" {
		return errors.ErrMetricNameMissing
	}

	if metricType != types.Counter && metricType != types.Gauge {
		return errors.ErrMetricTypeInvalid
	}

	return nil
}

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

func ValidateMetricID(id types.MetricID) error {
	if id.ID == "" {
		return errors.ErrMetricIDInvalid
	}

	if id.Type != types.Counter && id.Type != types.Gauge {
		return errors.ErrMetricTypeInvalid
	}

	return nil
}

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
