package core

import (
	"context"

	"github.com/glwbr/brisa/core/scraper"
)

// Scraper is the interface for portal-specific scrapers.
type Scraper interface {
	// GetCaptcha fetches a new captcha challenge.
	GetCaptcha(ctx context.Context) (*scraper.CaptchaChallenge, error)

	// SubmitWithCaptcha submits an access key with a pre-solved captcha.
	SubmitWithCaptcha(ctx context.Context, accessKey, captchaSolution string) (*ScrapingResult, error)
}

// ScrapingResult is the generic result of a scraping operation.
type ScrapingResult struct {
	// Receipt is the parsed invoice data.
	Receipt *Receipt
	// RawHTML holds the raw HTML responses by page name.
	RawHTML map[string][]byte
}

// ScraperConfig holds common configuration for scrapers.
type ScraperConfig struct {
	// CaptchaSolver is the solver to use for captcha challenges.
	CaptchaSolver scraper.CaptchaSolver
}

// NewScraper creates a scraper for the specified portal.
func NewScraper(portal Portal, cfg ScraperConfig) (Scraper, error) {
	switch portal {
	case PortalBA:
		opts := []BAScraperOption{}
		if cfg.CaptchaSolver != nil {
			opts = append(opts, WithCaptchaSolver(cfg.CaptchaSolver))
		}
		ba, err := NewBAScraper(opts...)
		if err != nil {
			return nil, err
		}
		return &baScraperAdapter{ba}, nil
	default:
		return nil, ErrUnsupportedPortal
	}
}

// baScraperAdapter adapts BAScraper to the Scraper interface.
type baScraperAdapter struct {
	s *BAScraper
}

func (a *baScraperAdapter) GetCaptcha(ctx context.Context) (*scraper.CaptchaChallenge, error) {
	return a.s.GetCaptcha(ctx)
}

func (a *baScraperAdapter) SubmitWithCaptcha(ctx context.Context, accessKey, captchaSolution string) (*ScrapingResult, error) {
	result, err := a.s.SubmitWithCaptcha(ctx, accessKey, captchaSolution)
	if err != nil {
		return nil, err
	}

	return &ScrapingResult{
		Receipt: result.Receipt,
		RawHTML: map[string][]byte{
			"danfe":    result.DanfeHTML,
			"nfe_tab":  result.NFETabHTML,
			"products": result.ProductsTabHTML,
		},
	}, nil
}
