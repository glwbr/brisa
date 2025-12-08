package core

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/glwbr/brisa/core/scraper"
)

// BAScraper handles the multi-step scraping flow for the BA SEFAZ NFC-e portal.
type BAScraper struct {
	client        *scraper.Client
	captchaSolver scraper.CaptchaSolver
	formState     *scraper.FormState
}

// BAScraperOption configures the BAScraper.
type BAScraperOption func(*BAScraper)

// WithCaptchaSolver sets the captcha solver for the scraper.
func WithCaptchaSolver(solver scraper.CaptchaSolver) BAScraperOption {
	return func(s *BAScraper) {
		s.captchaSolver = solver
	}
}

// NewBAScraper creates a new BA portal scraper.
func NewBAScraper(opts ...BAScraperOption) (*BAScraper, error) {
	// BA SEFAZ uses a certificate that may not be trusted by default
	client, err := scraper.NewClient(
		BABaseURL,
		scraper.WithInsecureSkipVerify(),
		scraper.WithHeader("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7"),
	)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	s := &BAScraper{
		client: client,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

// BAScrapingResult holds the result of a scraping session.
type BAScrapingResult struct {
	// Receipt is the parsed invoice data.
	Receipt *Receipt
	// NFETabHTML is the raw HTML of the NFe tab.
	NFETabHTML []byte
	// ProductsTabHTML is the raw HTML of the Products tab.
	ProductsTabHTML []byte
	// DanfeHTML is the raw HTML of the DANFE page.
	DanfeHTML []byte
}

// BAScrapingState tracks the current state of the scraping flow.
type BAScrapingState int

const (
	BAStateInitial BAScrapingState = iota
	BAStateAccessKeyPage
	BAStateCaptchaRequired
	BAStateDanfePage
	BAStateTabsPage
	BAStateComplete
	BAStateError
)

// String returns a human-readable state name.
func (s BAScrapingState) String() string {
	switch s {
	case BAStateInitial:
		return "initial"
	case BAStateAccessKeyPage:
		return "access_key_page"
	case BAStateCaptchaRequired:
		return "captcha_required"
	case BAStateDanfePage:
		return "danfe_page"
	case BAStateTabsPage:
		return "tabs_page"
	case BAStateComplete:
		return "complete"
	case BAStateError:
		return "error"
	default:
		return "unknown"
	}
}

// Common scraping errors.
var (
	ErrAccessKeyRequired  = errors.New("access key is required")
	ErrCaptchaSolverNil   = errors.New("captcha solver is required")
	ErrInvalidAccessKey   = errors.New("invalid access key format")
	ErrInvoiceNotFound    = errors.New("invoice not found")
	ErrUnexpectedResponse = errors.New("unexpected response from server")
	ErrSessionExpired     = errors.New("session expired, please retry")
)

// FetchByAccessKey performs the complete scraping flow for a given access key.
// It handles captcha solving, navigation, and tab parsing.
func (s *BAScraper) FetchByAccessKey(ctx context.Context, accessKey string) (*BAScrapingResult, error) {
	accessKey = normalizeAccessKey(accessKey)
	if !isValidAccessKey(accessKey) {
		return nil, ErrInvalidAccessKey
	}

	// Step 1: Load the access key page to get initial form state
	if err := s.loadAccessKeyPage(ctx); err != nil {
		return nil, fmt.Errorf("load access key page: %w", err)
	}

	// Step 2: Fetch and solve captcha
	captcha, err := s.fetchCaptcha(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch captcha: %w", err)
	}

	if s.captchaSolver == nil {
		return nil, ErrCaptchaSolverNil
	}

	solution, err := s.captchaSolver.Solve(ctx, captcha)
	if err != nil {
		return nil, fmt.Errorf("solve captcha: %w", err)
	}

	// Step 3: Submit access key + captcha, get DANFE page
	danfeHTML, err := s.submitAccessKey(ctx, accessKey, solution.Text)
	if err != nil {
		return nil, fmt.Errorf("submit access key: %w", err)
	}

	// Step 4: Navigate to tabs view
	tabsHTML, err := s.navigateToTabs(ctx, danfeHTML)
	if err != nil {
		return nil, fmt.Errorf("navigate to tabs: %w", err)
	}

	// Step 5: Parse NFe tab (already loaded)
	nfeTabHTML := tabsHTML

	// Step 6: Navigate to Products tab and parse
	productsTabHTML, err := s.loadTab(ctx, tabsHTML, BATabTypeProdutos)
	if err != nil {
		return nil, fmt.Errorf("load products tab: %w", err)
	}

	// Step 7: Parse the data
	result, err := s.parseResult(nfeTabHTML, productsTabHTML, danfeHTML)
	if err != nil {
		return nil, fmt.Errorf("parse result: %w", err)
	}

	return result, nil
}

// GetCaptcha fetches a new captcha challenge without submitting anything.
// Useful for UI flows that need to display the captcha first.
func (s *BAScraper) GetCaptcha(ctx context.Context) (*scraper.CaptchaChallenge, error) {
	// Ensure we have a fresh session
	if err := s.loadAccessKeyPage(ctx); err != nil {
		return nil, fmt.Errorf("load access key page: %w", err)
	}

	return s.fetchCaptcha(ctx)
}

// SubmitWithCaptcha submits an access key with a pre-solved captcha.
// Use this when captcha solving is handled externally (e.g., via UI).
func (s *BAScraper) SubmitWithCaptcha(ctx context.Context, accessKey, captchaSolution string) (*BAScrapingResult, error) {
	accessKey = normalizeAccessKey(accessKey)
	if !isValidAccessKey(accessKey) {
		return nil, ErrInvalidAccessKey
	}

	// Ensure we have form state
	if s.formState == nil || !s.formState.IsValid() {
		if err := s.loadAccessKeyPage(ctx); err != nil {
			return nil, fmt.Errorf("load access key page: %w", err)
		}
	}

	// Submit and continue the flow
	danfeHTML, err := s.submitAccessKey(ctx, accessKey, captchaSolution)
	if err != nil {
		return nil, fmt.Errorf("submit access key: %w", err)
	}

	tabsHTML, err := s.navigateToTabs(ctx, danfeHTML)
	if err != nil {
		return nil, fmt.Errorf("navigate to tabs: %w", err)
	}

	productsTabHTML, err := s.loadTab(ctx, tabsHTML, BATabTypeProdutos)
	if err != nil {
		return nil, fmt.Errorf("load products tab: %w", err)
	}

	return s.parseResult(tabsHTML, productsTabHTML, danfeHTML)
}

// loadAccessKeyPage loads the initial page and extracts form state.
func (s *BAScraper) loadAccessKeyPage(ctx context.Context) error {
	resp, err := s.client.Get(ctx, BAAccessKeyPage, nil)
	if err != nil {
		return err
	}

	body, err := resp.Body()
	if err != nil {
		return err
	}

	state, err := scraper.ParseFormState(body)
	if err != nil {
		return fmt.Errorf("parse form state: %w", err)
	}

	s.formState = state
	return nil
}

// fetchCaptcha retrieves a captcha image from the portal.
func (s *BAScraper) fetchCaptcha(ctx context.Context) (*scraper.CaptchaChallenge, error) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	imageData, contentType, err := s.client.GetImage(ctx, BACaptchaEndpoint, &scraper.RequestConfig{
		Params: url.Values{
			"t": {strconv.FormatInt(timestamp, 10)},
		},
		Referer: s.client.BaseURL() + BAAccessKeyPage,
	})
	if err != nil {
		return nil, scraper.NewCaptchaFetchError(err)
	}

	return &scraper.CaptchaChallenge{
		ID:          strconv.FormatInt(timestamp, 10),
		Image:       imageData,
		ContentType: contentType,
	}, nil
}

// submitAccessKey submits the access key and captcha, returning the DANFE page HTML.
func (s *BAScraper) submitAccessKey(ctx context.Context, accessKey, captchaSolution string) ([]byte, error) {
	if s.formState == nil {
		return nil, errors.New("form state not initialized")
	}

	formData := scraper.NewFormBuilder(s.formState).
		Set(BAFieldAccessKey, accessKey).
		Set(BAFieldCaptcha, captchaSolution).
		Set(BAFieldSubmit, "Consultar").
		Build()

	resp, err := s.client.PostForm(ctx, BAAccessKeyPage, toURLValues(formData), &scraper.RequestConfig{
		Referer: s.client.BaseURL() + BAAccessKeyPage,
	})
	if err != nil {
		return nil, err
	}

	body, err := resp.Body()
	if err != nil {
		return nil, err
	}

	// Check for error responses
	if err := s.checkForErrors(body); err != nil {
		return nil, err
	}

	// Update form state from response
	if newState, err := scraper.ParseFormState(body); err == nil && newState.IsValid() {
		s.formState = newState
	}

	return body, nil
}

// navigateToTabs clicks "Visualizar em Abas" to get the tabbed view.
func (s *BAScraper) navigateToTabs(ctx context.Context, danfeHTML []byte) ([]byte, error) {
	// Parse form state from DANFE page
	state, err := scraper.ParseFormState(danfeHTML)
	if err != nil {
		return nil, fmt.Errorf("parse danfe form state: %w", err)
	}
	s.formState = state

	formData := scraper.NewFormBuilder(s.formState).
		Set(BAFieldViewTabs, "Visualizar em Abas").
		Set(BAFieldOriginCall, "").
		Build()

	resp, err := s.client.PostForm(ctx, BADanfePage, toURLValues(formData), &scraper.RequestConfig{
		Referer: s.client.BaseURL() + BADanfePage,
	})
	if err != nil {
		return nil, err
	}

	body, err := resp.Body()
	if err != nil {
		return nil, err
	}

	// Update form state
	if newState, err := scraper.ParseFormState(body); err == nil && newState.IsValid() {
		s.formState = newState
	}

	return body, nil
}

// loadTab navigates to a specific tab and returns the HTML.
func (s *BAScraper) loadTab(ctx context.Context, currentHTML []byte, tab BATab) ([]byte, error) {
	// Parse form state from current page
	state, err := scraper.ParseFormState(currentHTML)
	if err != nil {
		return nil, fmt.Errorf("parse tab form state: %w", err)
	}
	s.formState = state

	buttonName := BATabButtonName(tab)
	if buttonName == "" {
		return nil, fmt.Errorf("unknown tab: %s", tab)
	}

	// For image buttons, we need to send .x and .y coordinates
	formData := scraper.NewFormBuilder(s.formState).
		Set(buttonName+".x", "10").
		Set(buttonName+".y", "10").
		Set(BAFieldOriginCall, "").
		Build()

	resp, err := s.client.PostForm(ctx, BATabsPage, toURLValues(formData), &scraper.RequestConfig{
		Referer: s.client.BaseURL() + BATabsPage,
	})
	if err != nil {
		return nil, err
	}

	body, err := resp.Body()
	if err != nil {
		return nil, err
	}

	// Update form state
	if newState, err := scraper.ParseFormState(body); err == nil && newState.IsValid() {
		s.formState = newState
	}

	return body, nil
}

// checkForErrors examines the response HTML for error messages.
func (s *BAScraper) checkForErrors(html []byte) error {
	htmlStr := string(html)

	// Check for common error patterns
	errorPatterns := map[string]error{
		"Chave de Acesso inválida":                             ErrInvalidAccessKey,
		"chave de acesso inválida":                             ErrInvalidAccessKey,
		"NFC-e não encontrada":                                 ErrInvoiceNotFound,
		"Captcha inválido":                                     scraper.ErrCaptchaInvalid,
		"código de segurança inválido":                         scraper.ErrCaptchaInvalid,
		"código impresso ao lado":                              scraper.ErrCaptchaInvalid, // Still on captcha page
		"sessão expirou":                                       ErrSessionExpired,
		"Sessão expirada":                                      ErrSessionExpired,
		"Object reference not set to an instance of an object": ErrUnexpectedResponse,
		"Ocorreu um erro":                                      ErrUnexpectedResponse,
	}

	for pattern, err := range errorPatterns {
		if strings.Contains(strings.ToLower(htmlStr), strings.ToLower(pattern)) {
			return fmt.Errorf("%w: %s", err, pattern)
		}
	}

	return nil
}

// parseResult combines the parsed data from multiple tabs.
func (s *BAScraper) parseResult(nfeTabHTML, productsTabHTML, danfeHTML []byte) (*BAScrapingResult, error) {
	// Parse NFe tab for receipt header data
	receipt, err := ParseBAFromAbasNFe(nfeTabHTML)
	if err != nil {
		return nil, fmt.Errorf("parse nfe tab: %w", err)
	}

	// Parse Products tab for items
	items, err := ParseBAFromAbasProdutosServicos(productsTabHTML)
	if err != nil {
		// Items are optional - some receipts may not have the products tab
		items = []Item{}
	}

	receipt.Items = items

	return &BAScrapingResult{
		Receipt:         receipt,
		NFETabHTML:      nfeTabHTML,
		ProductsTabHTML: productsTabHTML,
		DanfeHTML:       danfeHTML,
	}, nil
}

// normalizeAccessKey removes spaces and formatting from an access key.
func normalizeAccessKey(key string) string {
	var b strings.Builder
	for _, r := range key {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// isValidAccessKey checks if an access key has the correct format.
// Brazilian NFC-e access keys are 44 digits.
func isValidAccessKey(key string) bool {
	return len(key) == 44
}

// toURLValues converts a map to url.Values.
func toURLValues(m map[string]string) url.Values {
	v := url.Values{}
	for k, val := range m {
		v.Set(k, val)
	}
	return v
}
