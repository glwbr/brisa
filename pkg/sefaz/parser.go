package sefaz

import (
	"context"

	"github.com/glwbr/brisa/internal/client"
	"github.com/glwbr/brisa/pkg/invoice"
	"github.com/glwbr/brisa/pkg/logger"
)

// Parser defines the interface that all SEFAZ parsers must implement
type Parser interface {
	// Initialize sets up the parser with initial state
	Initialize(ctx context.Context) error

	// GetCaptcha retrieves a captcha image
	GetCaptcha(ctx context.Context) (CaptchaImage, error)

	// Submit submits a form with key and captcha resolution
	Submit(ctx context.Context, key, solution string) error

	// FetchInvoice retrieves an invoice by its identification
	FetchInvoice(ctx context.Context) (*invoice.Invoice, error)

	// FetchDetailedInvoice retrieves a detailed invoice view // maybe invoice.DetailedInvoice
	FetchDetailedInvoice(ctx context.Context, key string) (*invoice.Invoice, error)
}

// ParserOption represents a configuration option for parsers
type ParserOption func(*ParserOptions)

// ParserOptions holds configuration for parsers
type ParserOptions struct {
	Client    *client.Client
	Logger    logger.Logger
	CacheDir  string
	UserAgent string
	// Add more options as needed
}

// CaptchaImage represents a captcha challenge
type CaptchaImage struct {
	Data      []byte
	MimeType  string
	Challenge string
}
