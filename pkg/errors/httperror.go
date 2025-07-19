package errors

import (
	"fmt"
	"net/http"
)

// HTTPError represents an error that occurred during an HTTP operation.
// It captures HTTP-specific context including status code, URL, method,
// and optionally the response body and underlying error.
type HTTPError struct {
	StatusCode int
	URL        string
	Method     string
	Message    string
	RawError   error
	Response   any
}

// Error implements the error interface, providing a formatted string
// representation of the HTTP error including method, URL, status code,
// and any additional messages or underlying errors.
func (e *HTTPError) Error() string {
	msg := fmt.Sprintf("%s %s: %d", e.Method, e.URL, e.StatusCode)
	if e.Message != "" {
		msg += " - " + e.Message
	}
	if e.RawError != nil {
		msg += fmt.Sprintf(" (caused by: %s)", e.RawError)
	}
	return msg
}

// Unwrap returns the underlying error that caused this HTTPError,
// allowing errors.Is and errors.As to work with the error chain.
func (e *HTTPError) Unwrap() error {
	return e.RawError
}

// NewHTTPError creates a new HTTPError from an HTTP response.
// If resp is nil, creates an error with status code 0 (unknown).
// The message parameter provides additional context about the error.
func NewHTTPError(resp *http.Response, err error, message string) *HTTPError {
	if resp == nil {
		return &HTTPError{
			StatusCode: 0,
			Message:    message,
			RawError:   err,
		}
	}

	return &HTTPError{
		StatusCode: resp.StatusCode,
		URL:        resp.Request.URL.String(),
		Method:     resp.Request.Method,
		Message:    message,
		RawError:   err,
	}
}

// IsNotFound reports whether the error represents an HTTP 404 Not Found response.
// Returns true if the error is an HTTPError with StatusCode 404,
// or wraps such an error in its chain.
func IsNotFound(err error) bool {
	var httpErr *HTTPError
	if As(err, &httpErr) {
		return httpErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsServerError reports whether the error represents an HTTP 5xx server error.
// Returns true if the error is an HTTPError with StatusCode between 500-599,
// or wraps such an error in its chain.
func IsServerError(err error) bool {
	var httpErr *HTTPError
	if As(err, &httpErr) {
		return httpErr.StatusCode >= http.StatusInternalServerError && httpErr.StatusCode < 600
	}
	return false
}

// IsClientError reports whether the error represents an HTTP 4xx client error.
// Returns true if the error is an HTTPError with StatusCode between 400-499,
// or wraps such an error in its chain.
func IsClientError(err error) bool {
	var httpErr *HTTPError
	if As(err, &httpErr) {
		return httpErr.StatusCode >= http.StatusBadRequest && httpErr.StatusCode < http.StatusInternalServerError
	}
	return false
}
