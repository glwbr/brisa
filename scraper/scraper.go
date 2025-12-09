// Package scraper defines interfaces and common types for NFC-e portal scraping.
package scraper

import (
	"context"
	"errors"

	"github.com/glwbr/brisa/invoice"
)

// TODO: normalize the errors and its messages for the API
var (
	ErrAccessKeyRequired  = errors.New("access key is required")
	ErrInvalidAccessKey   = errors.New("invalid access key format")
	ErrInvoiceNotFound    = errors.New("invoice not found")
	ErrSessionExpired     = errors.New("session expired")
	ErrUnexpectedResponse = errors.New("unexpected server response")
)

// Scraper fetches invoice data from a portal.
type Scraper interface {
	GetCaptcha(ctx context.Context) (*CaptchaChallenge, error)
	SubmitWithCaptcha(ctx context.Context, accessKey, captchaSolution string) (*Result, error)
}

// Result holds the outcome of a scraping operation.
type Result struct {
	Receipt *invoice.Receipt
	RawHTML map[string][]byte
}
