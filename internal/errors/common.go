// Package errors defines application-specific error variables
// related to metric processing.
package errors

import "errors"

var (
	// ErrInternalServerError indicates that an unexpected server-side error occurred.
	ErrInternalServerError = errors.New("internal server error")
)
