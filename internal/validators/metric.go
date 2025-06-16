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
	case string(types.Counter):
		_, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.ErrMetricValueInvalid
		}
	case string(types.Gauge):
		_, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors.ErrMetricValueInvalid
		}
	}

	return nil
}
