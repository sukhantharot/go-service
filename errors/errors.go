package errors

import (
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// New creates a new AppError
func New(code int, message string, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common errors
var (
	ErrUnauthorized = New(http.StatusUnauthorized, "Unauthorized", "")
	ErrForbidden    = New(http.StatusForbidden, "Forbidden", "")
	ErrNotFound     = New(http.StatusNotFound, "Not Found", "")
	ErrBadRequest   = New(http.StatusBadRequest, "Bad Request", "")
	ErrInternal     = New(http.StatusInternalServerError, "Internal Server Error", "")
)

// Database errors
var (
	ErrDatabaseConnection = New(http.StatusInternalServerError, "Database Connection Error", "")
	ErrRecordNotFound     = New(http.StatusNotFound, "Record Not Found", "")
	ErrDuplicateEntry     = New(http.StatusBadRequest, "Duplicate Entry", "")
)

// Authentication errors
var (
	ErrInvalidCredentials = New(http.StatusUnauthorized, "Invalid Credentials", "")
	ErrTokenExpired       = New(http.StatusUnauthorized, "Token Expired", "")
	ErrTokenInvalid       = New(http.StatusUnauthorized, "Invalid Token", "")
)

// Validation errors
var (
	ErrValidation = New(http.StatusBadRequest, "Validation Error", "")
) 