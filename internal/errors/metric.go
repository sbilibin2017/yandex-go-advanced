package errors

import "errors"

var (
	ErrMetricNameMissing  = errors.New("metric name is required")
	ErrMetricTypeInvalid  = errors.New("invalid metric type")
	ErrMetricNotFound     = errors.New("metric not found")
	ErrMetricIDInvalid    = errors.New("invalid metric id")
	ErrMetricValueInvalid = errors.New("invalid metric value")
	ErrMetricDeltaInvalid = errors.New("invalid metric delta")
)
