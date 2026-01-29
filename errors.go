package filemaker

import (
	"errors"
	"fmt"
	"net/http"
)

// FileMakerError represents an error returned by the FileMaker Data API.
// It includes the error code, message, and HTTP status code.
type FileMakerError struct {
	Code       string // FileMaker error code (e.g., "401", "952")
	Message    string // Error message from FileMaker
	HTTPStatus int    // HTTP status code
	Err        error  // Underlying error if any
}

// Error implements the error interface.
func (e *FileMakerError) Error() string {
	if e.Code != "" && e.Message != "" {
		return fmt.Sprintf("FileMaker error %s: %s (HTTP %d)", e.Code, e.Message, e.HTTPStatus)
	}
	if e.Message != "" {
		return fmt.Sprintf("FileMaker error: %s (HTTP %d)", e.Message, e.HTTPStatus)
	}
	return fmt.Sprintf("FileMaker error (HTTP %d)", e.HTTPStatus)
}

// Unwrap returns the underlying error.
func (e *FileMakerError) Unwrap() error {
	return e.Err
}

// ValidationError represents an input validation error.
type ValidationError struct {
	Field   string // Field name that failed validation
	Message string // Validation error message
	Err     error  // Underlying error if any
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// Unwrap returns the underlying error.
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// AuthenticationError represents an authentication failure.
type AuthenticationError struct {
	Message string
	Err     error
}

// Error implements the error interface.
func (e *AuthenticationError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("authentication error: %s", e.Message)
	}
	return "authentication error"
}

// Unwrap returns the underlying error.
func (e *AuthenticationError) Unwrap() error {
	return e.Err
}

// NetworkError represents a network-level error.
type NetworkError struct {
	Message string
	Err     error
}

// Error implements the error interface.
func (e *NetworkError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("network error: %s", e.Message)
	}
	return "network error"
}

// Unwrap returns the underlying error.
func (e *NetworkError) Unwrap() error {
	return e.Err
}

// TimeoutError represents a timeout error.
type TimeoutError struct {
	Message string
	Err     error
}

// Error implements the error interface.
func (e *TimeoutError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("timeout error: %s", e.Message)
	}
	return "timeout error"
}

// Unwrap returns the underlying error.
func (e *TimeoutError) Unwrap() error {
	return e.Err
}

// Error checking helpers

// IsFileMakerError checks if an error is a FileMakerError.
func IsFileMakerError(err error) bool {
	var fmErr *FileMakerError
	return errors.As(err, &fmErr)
}

// IsValidationError checks if an error is a ValidationError.
func IsValidationError(err error) bool {
	var valErr *ValidationError
	return errors.As(err, &valErr)
}

// IsAuthenticationError checks if an error is an AuthenticationError.
func IsAuthenticationError(err error) bool {
	var authErr *AuthenticationError
	return errors.As(err, &authErr)
}

// IsNetworkError checks if an error is a NetworkError.
func IsNetworkError(err error) bool {
	var netErr *NetworkError
	return errors.As(err, &netErr)
}

// IsTimeoutError checks if an error is a TimeoutError.
func IsTimeoutError(err error) bool {
	var timeoutErr *TimeoutError
	return errors.As(err, &timeoutErr)
}

// IsRetryable determines if an error is retryable based on FileMaker error codes
// and HTTP status codes.
//
// Retryable errors include:
// - HTTP 429 (Too Many Requests)
// - HTTP 500, 502, 503, 504 (Server errors)
// - FileMaker error 952 (Host unavailable)
// - FileMaker error 953 (Too many files open)
// - Network errors
// - Timeout errors
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check FileMaker errors
	var fmErr *FileMakerError
	if errors.As(err, &fmErr) {
		// Retryable FileMaker error codes
		switch fmErr.Code {
		case "952": // Host unavailable
			return true
		case "953": // Too many files open
			return true
		}

		// Retryable HTTP status codes
		switch fmErr.HTTPStatus {
		case http.StatusTooManyRequests:
			return true
		case http.StatusInternalServerError:
			return true
		case http.StatusBadGateway:
			return true
		case http.StatusServiceUnavailable:
			return true
		case http.StatusGatewayTimeout:
			return true
		}

		return false
	}

	// Network and timeout errors are retryable
	if IsNetworkError(err) || IsTimeoutError(err) {
		return true
	}

	return false
}

// IsAuthError checks if the error is related to authentication.
// This includes both AuthenticationError and FileMaker auth-related error codes.
func IsAuthError(err error) bool {
	if IsAuthenticationError(err) {
		return true
	}

	var fmErr *FileMakerError
	if errors.As(err, &fmErr) {
		// FileMaker authentication error codes
		switch fmErr.Code {
		case "212": // Invalid username or password
			return true
		case "214": // Account name does not have access privileges
			return true
		case "952": // Invalid FileMaker Data API token
			return true
		}

		// HTTP auth status codes
		if fmErr.HTTPStatus == http.StatusUnauthorized || fmErr.HTTPStatus == http.StatusForbidden {
			return true
		}
	}

	return false
}

// ParseFileMakerError creates a FileMakerError from a ResponseData.
// It extracts the first error code and message from the Messages array.
func ParseFileMakerError(response *ResponseData, httpStatus int) error {
	if response == nil {
		return &FileMakerError{
			HTTPStatus: httpStatus,
			Message:    "no response data",
		}
	}

	if len(response.Messages) == 0 {
		return &FileMakerError{
			HTTPStatus: httpStatus,
			Message:    "unknown error",
		}
	}

	// Use the first message
	msg := response.Messages[0]

	// Error code "0" means success in FileMaker
	if msg.Code == "0" {
		return nil
	}

	return &FileMakerError{
		Code:       msg.Code,
		Message:    msg.Message,
		HTTPStatus: httpStatus,
	}
}
