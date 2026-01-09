package http

import "errors"

// HTTP adapter errors.
var (
	// ErrEmptyRequestBody is returned when the request body is empty.
	ErrEmptyRequestBody = errors.New("request body is empty")

	// ErrInvalidUUID is returned when a UUID parameter is invalid.
	ErrInvalidUUID = errors.New("invalid UUID format")

	// ErrUnauthorized is returned when authentication fails.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when the user lacks permission.
	ErrForbidden = errors.New("forbidden")
)
