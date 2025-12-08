package scraper

import "errors"

// Common scraping errors.
var (
	// ErrTimeout indicates a request timeout.
	ErrTimeout = errors.New("request timeout")

	// ErrNetworkError indicates a network-level error.
	ErrNetworkError = errors.New("network error")

	// ErrInvalidResponse indicates the server returned an unexpected response.
	ErrInvalidResponse = errors.New("invalid server response")

	// ErrRateLimited indicates the request was rate limited.
	ErrRateLimited = errors.New("rate limited")

	// ErrSessionInvalid indicates the session is invalid or expired.
	ErrSessionInvalid = errors.New("session invalid")
)

// ScrapingError wraps errors with additional context.
type ScrapingError struct {
	Step    string // The step where the error occurred
	Message string
	Cause   error
}

func (e *ScrapingError) Error() string {
	if e.Cause != nil {
		return e.Step + ": " + e.Message + ": " + e.Cause.Error()
	}
	return e.Step + ": " + e.Message
}

func (e *ScrapingError) Unwrap() error {
	return e.Cause
}

// NewScrapingError creates a new scraping error.
func NewScrapingError(step, message string, cause error) *ScrapingError {
	return &ScrapingError{
		Step:    step,
		Message: message,
		Cause:   cause,
	}
}

// IsRetryable returns true if the error is potentially recoverable by retrying.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific retryable errors
	if errors.Is(err, ErrTimeout) ||
		errors.Is(err, ErrNetworkError) ||
		errors.Is(err, ErrRateLimited) ||
		errors.Is(err, ErrCaptchaExpired) {
		return true
	}

	return false
}
