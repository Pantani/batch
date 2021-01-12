package api

import (
	"github.com/Pantani/errors"
)

var (
	// ErrInvalidID signals that the requested id is invalid
	ErrInvalidID = errors.E("invalid id")
)

type (
	errResponse struct {
		Error errorDetails `json:"error"`
	}
	errorDetails struct {
		Message string `json:"message"`
	}
)

// errorResponse returns the error response object.
func errorResponse(err error) errResponse {
	var message string
	if err != nil {
		message = err.Error()
	}
	return errResponse{Error: errorDetails{
		Message: message,
	}}
}
