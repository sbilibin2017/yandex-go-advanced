package errors

import "errors"

var (
	ErrMetricNameMissing  = errors.New("metric name is required")
	ErrMetricTypeInvalid  = errors.New("invalid metric type")
	ErrMetricValueInvalid = errors.New("invalid metric value")
	ErrMetricNotFound     = errors.New("metric not found")
)
