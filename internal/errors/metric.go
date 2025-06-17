// Package errors defines common application-level error variables
// used for validating and processing metrics.
package errors

import "errors"

var (
	// ErrMetricNameMissing indicates that a required metric name is missing.
	ErrMetricNameMissing = errors.New("metric name is required")

	// ErrMetricTypeInvalid indicates that the provided metric type is not supported or invalid.
	ErrMetricTypeInvalid = errors.New("invalid metric type")

	// ErrMetricNotFound indicates that the requested metric could not be found.
	ErrMetricNotFound = errors.New("metric not found")

	// ErrMetricIDInvalid indicates that the provided metric ID is invalid or malformed.
	ErrMetricIDInvalid = errors.New("invalid metric id")

	// ErrMetricValueInvalid indicates that a provided metric value is invalid or cannot be processed.
	ErrMetricValueInvalid = errors.New("invalid metric value")

	// ErrMetricDeltaInvalid indicates that a provided metric delta is invalid or cannot be processed.
	ErrMetricDeltaInvalid = errors.New("invalid metric delta")
)
