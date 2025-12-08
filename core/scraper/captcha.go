package scraper

import (
	"context"
)

// CaptchaChallenge represents a captcha that needs to be solved.
type CaptchaChallenge struct {
	// ID is a unique identifier for this challenge (optional, portal-specific).
	ID string
	// Image is the captcha image bytes.
	Image []byte
	// ContentType is the MIME type of the image (e.g., "image/png").
	ContentType string
	// Metadata holds portal-specific additional data.
	Metadata map[string]string
}

// CaptchaSolution holds the user's captcha answer.
type CaptchaSolution struct {
	// Text is the user's solution text.
	Text string
	// ChallengeID references the original challenge (if applicable).
	ChallengeID string
}

// CaptchaSolver is an interface for captcha solving strategies.
// This allows different implementations: manual, OCR, or third-party services.
type CaptchaSolver interface {
	// Solve presents the captcha to the user/service and returns the solution.
	Solve(ctx context.Context, challenge *CaptchaChallenge) (*CaptchaSolution, error)
}

// CaptchaProvider fetches captcha challenges from a portal.
type CaptchaProvider interface {
	// FetchCaptcha retrieves a new captcha challenge.
	FetchCaptcha(ctx context.Context) (*CaptchaChallenge, error)
	// RefreshCaptcha gets a new captcha (e.g., if user can't read the current one).
	RefreshCaptcha(ctx context.Context) (*CaptchaChallenge, error)
}

// ManualSolver is a placeholder for manual captcha solving.
// In practice, this would integrate with UI callbacks.
type ManualSolver struct {
	// PromptFunc is called to display the captcha and get user input.
	PromptFunc func(ctx context.Context, challenge *CaptchaChallenge) (string, error)
}

// Solve implements CaptchaSolver for manual solving.
func (s *ManualSolver) Solve(ctx context.Context, challenge *CaptchaChallenge) (*CaptchaSolution, error) {
	if s.PromptFunc == nil {
		return nil, ErrNoCaptchaSolver
	}

	text, err := s.PromptFunc(ctx, challenge)
	if err != nil {
		return nil, err
	}

	return &CaptchaSolution{
		Text:        text,
		ChallengeID: challenge.ID,
	}, nil
}

// CaptchaCallback is a function type for receiving captcha challenges.
// This enables event-driven UI updates without tight coupling.
type CaptchaCallback func(challenge *CaptchaChallenge) (*CaptchaSolution, error)

// CallbackSolver wraps a callback function as a CaptchaSolver.
type CallbackSolver struct {
	Callback CaptchaCallback
}

// Solve implements CaptchaSolver using a callback.
func (s *CallbackSolver) Solve(ctx context.Context, challenge *CaptchaChallenge) (*CaptchaSolution, error) {
	if s.Callback == nil {
		return nil, ErrNoCaptchaSolver
	}
	return s.Callback(challenge)
}

// CaptchaError represents captcha-related errors.
type CaptchaError struct {
	Code    string
	Message string
	Cause   error
}

func (e *CaptchaError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *CaptchaError) Unwrap() error {
	return e.Cause
}

// Common captcha errors.
var (
	ErrNoCaptchaSolver = &CaptchaError{
		Code:    "NO_SOLVER",
		Message: "no captcha solver configured",
	}
	ErrCaptchaExpired = &CaptchaError{
		Code:    "EXPIRED",
		Message: "captcha has expired",
	}
	ErrCaptchaInvalid = &CaptchaError{
		Code:    "INVALID",
		Message: "invalid captcha solution",
	}
	ErrCaptchaFetchFailed = &CaptchaError{
		Code:    "FETCH_FAILED",
		Message: "failed to fetch captcha",
	}
)

// NewCaptchaFetchError creates a captcha fetch error with a cause.
func NewCaptchaFetchError(cause error) *CaptchaError {
	return &CaptchaError{
		Code:    "FETCH_FAILED",
		Message: "failed to fetch captcha",
		Cause:   cause,
	}
}

// NewCaptchaInvalidError creates an invalid captcha error with details.
func NewCaptchaInvalidError(details string) *CaptchaError {
	return &CaptchaError{
		Code:    "INVALID",
		Message: "invalid captcha solution: " + details,
	}
}
