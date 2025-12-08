package scraper

import (
	"context"
	"errors"
)

var (
	ErrNoCaptchaSolver = errors.New("no captcha solver configured")
	ErrCaptchaExpired  = errors.New("captcha expired")
	ErrCaptchaInvalid  = errors.New("invalid captcha solution")
)

type CaptchaChallenge struct {
	ID          string
	Image       []byte
	ContentType string
	Metadata    map[string]string
}

type CaptchaSolution struct {
	Text        string
	ChallengeID string
}

// CaptchaSolver resolves captcha challenges.
type CaptchaSolver interface {
	Solve(ctx context.Context, challenge *CaptchaChallenge) (*CaptchaSolution, error)
}

// ManualSolver delegates captcha solving to a callback function.
type ManualSolver struct {
	PromptFunc func(ctx context.Context, challenge *CaptchaChallenge) (string, error)
}

func (s *ManualSolver) Solve(ctx context.Context, challenge *CaptchaChallenge) (*CaptchaSolution, error) {
	if s.PromptFunc == nil {
		return nil, ErrNoCaptchaSolver
	}
	text, err := s.PromptFunc(ctx, challenge)
	if err != nil {
		return nil, err
	}
	return &CaptchaSolution{Text: text, ChallengeID: challenge.ID}, nil
}
