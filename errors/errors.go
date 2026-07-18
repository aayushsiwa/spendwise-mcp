package errors

import "fmt"

type AppError struct {
	Type       string         `json:"type"`
	Message    string         `json:"message"`
	Details    map[string]any `json:"details,omitempty"`
	StatusCode int            `json:"-"`
	Err        error          `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(errType, message string, statusCode int, err error) *AppError {
	return &AppError{Type: errType, Message: message, StatusCode: statusCode, Err: err}
}

func (e *AppError) WithDetails(details map[string]any) *AppError {
	e.Details = details
	return e
}

func NewValidation(message string, details map[string]any) *AppError {
	return New("validation_error", message, 400, nil).WithDetails(details)
}

func NewUnauthorized(message string, err error) *AppError {
	return New("unauthorized", message, 401, err)
}

func NewForbidden(message string, err error) *AppError {
	return New("forbidden", message, 403, err)
}

func NewNotFound(message string, err error) *AppError {
	return New("not_found", message, 404, err)
}

func NewConflict(message string, err error) *AppError {
	return New("conflict", message, 409, err)
}

func NewRateLimited(message string, err error) *AppError {
	return New("rate_limited", message, 429, err)
}

func NewBackend(message string, err error) *AppError {
	return New("backend_error", message, 502, err)
}

func NewInternal(message string, err error) *AppError {
	return New("internal_error", message, 500, err)
}
