package ba

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/glwbr/brisa/internal/http"
	"github.com/glwbr/brisa/scraper"
)

type Scraper struct {
	client        *http.Client
	captchaSolver scraper.CaptchaSolver
	formState     *scraper.FormState
}

type Option func(*Scraper)

func WithCaptchaSolver(solver scraper.CaptchaSolver) Option {
	return func(s *Scraper) { s.captchaSolver = solver }
}

func New(opts ...Option) (*Scraper, error) {
	client, err := http.New(BaseURL, http.WithInsecureSkipVerify())
	if err != nil {
		return nil, err
	}
	s := &Scraper{client: client}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Scraper) GetCaptcha(ctx context.Context) (*scraper.CaptchaChallenge, error) {
	if err := s.loadAccessKeyPage(ctx); err != nil {
		return nil, err
	}
	return s.fetchCaptcha(ctx)
}

func (s *Scraper) SubmitWithCaptcha(ctx context.Context, accessKey, captchaSolution string) (*scraper.Result, error) {
	accessKey = normalizeAccessKey(accessKey)
	if len(accessKey) != 44 {
		return nil, scraper.ErrInvalidAccessKey
	}

	if s.formState == nil || !s.formState.IsValid() {
		if err := s.loadAccessKeyPage(ctx); err != nil {
			return nil, err
		}
	}

	danfeHTML, err := s.submitAccessKey(ctx, accessKey, captchaSolution)
	if err != nil {
		return nil, err
	}

	tabsHTML, err := s.navigateToTabs(ctx, danfeHTML)
	if err != nil {
		return nil, err
	}

	productsHTML, err := s.loadTab(ctx, tabsHTML, TabProdutos)
	if err != nil {
		return nil, err
	}

	receipt, err := ParseNFeTab(tabsHTML)
	if err != nil {
		return nil, fmt.Errorf("parse nfe tab: %w", err)
	}

	items, err := ParseProductsTab(productsHTML)
	if err == nil {
		receipt.Items = items
	}

	return &scraper.Result{
		Receipt: receipt,
		RawHTML: map[string][]byte{
			"danfe":    danfeHTML,
			"nfe_tab":  tabsHTML,
			"products": productsHTML,
		},
	}, nil
}

func (s *Scraper) FetchByAccessKey(ctx context.Context, accessKey string) (*scraper.Result, error) {
	if s.captchaSolver == nil {
		return nil, scraper.ErrNoCaptchaSolver
	}

	for {
		challenge, err := s.GetCaptcha(ctx)
		if err != nil {
			return nil, err
		}

		solution, err := s.captchaSolver.Solve(ctx, challenge)
		if err != nil {
			return nil, err
		}

		result, err := s.SubmitWithCaptcha(ctx, accessKey, solution.Text)
		if errors.Is(err, scraper.ErrCaptchaInvalid) {
			fmt.Println("Invalid captcha, retrying...")
			continue
		}
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (s *Scraper) loadAccessKeyPage(ctx context.Context) error {
	resp, err := s.client.Get(ctx, AccessKeyPage, nil)
	if err != nil {
		return err
	}
	body, err := resp.Body()
	if err != nil {
		return err
	}
	state, err := scraper.ParseFormState(body)
	if err != nil {
		return err
	}
	s.formState = state
	return nil
}

func (s *Scraper) fetchCaptcha(ctx context.Context) (*scraper.CaptchaChallenge, error) {
	ts := time.Now().UnixMilli()
	image, contentType, err := s.client.GetImage(ctx, CaptchaEndpoint, &http.RequestConfig{
		Params:  url.Values{"t": {strconv.FormatInt(ts, 10)}},
		Referer: s.client.BaseURL() + AccessKeyPage,
	})
	if err != nil {
		return nil, err
	}
	return &scraper.CaptchaChallenge{
		ID:          strconv.FormatInt(ts, 10),
		Image:       image,
		ContentType: contentType,
	}, nil
}

func (s *Scraper) submitAccessKey(ctx context.Context, accessKey, captcha string) ([]byte, error) {
	if s.formState == nil {
		return nil, errors.New("form state not initialized")
	}

	form := scraper.NewFormBuilder(s.formState).
		Set(FieldAccessKey, accessKey).
		Set(FieldCaptcha, captcha).
		Set(FieldSubmit, "Consultar").
		Build()

	resp, err := s.client.PostForm(ctx, AccessKeyPage, toURLValues(form), &http.RequestConfig{
		Referer: s.client.BaseURL() + AccessKeyPage,
	})
	if err != nil {
		return nil, err
	}

	body, err := resp.Body()
	if err != nil {
		return nil, err
	}

	if err := checkForErrors(body); err != nil {
		return nil, err
	}

	if state, err := scraper.ParseFormState(body); err == nil && state.IsValid() {
		s.formState = state
	}
	return body, nil
}

func (s *Scraper) navigateToTabs(ctx context.Context, danfeHTML []byte) ([]byte, error) {
	state, err := scraper.ParseFormState(danfeHTML)
	if err != nil {
		return nil, err
	}
	s.formState = state

	form := scraper.NewFormBuilder(s.formState).
		Set(FieldViewTabs, "Visualizar em Abas").
		Build()

	resp, err := s.client.PostForm(ctx, DanfePage, toURLValues(form), &http.RequestConfig{
		Referer: s.client.BaseURL() + DanfePage,
	})
	if err != nil {
		return nil, err
	}

	body, err := resp.Body()
	if err != nil {
		return nil, err
	}

	if state, err := scraper.ParseFormState(body); err == nil && state.IsValid() {
		s.formState = state
	}
	return body, nil
}

func (s *Scraper) loadTab(ctx context.Context, currentHTML []byte, tab Tab) ([]byte, error) {
	state, err := scraper.ParseFormState(currentHTML)
	if err != nil {
		return nil, err
	}
	s.formState = state

	btn := tab.ButtonName()
	if btn == "" {
		return nil, fmt.Errorf("unknown tab: %s", tab)
	}

	form := scraper.NewFormBuilder(s.formState).
		Set(btn+".x", "10").
		Set(btn+".y", "10").
		Build()

	resp, err := s.client.PostForm(ctx, TabsPage, toURLValues(form), &http.RequestConfig{
		Referer: s.client.BaseURL() + TabsPage,
	})
	if err != nil {
		return nil, err
	}

	body, err := resp.Body()
	if err != nil {
		return nil, err
	}

	if state, err := scraper.ParseFormState(body); err == nil && state.IsValid() {
		s.formState = state
	}
	return body, nil
}

func checkForErrors(html []byte) error {
	s := strings.ToLower(string(html))
	patterns := map[string]error{
		"chave de acesso inválida":         scraper.ErrInvalidAccessKey,
		"nfc-e não encontrada":             scraper.ErrInvoiceNotFound,
		"captcha inválido":                 scraper.ErrCaptchaInvalid,
		"código incorreto tente novamente": scraper.ErrCaptchaInvalid,
		"sessão expirada":                  scraper.ErrSessionExpired,
		"ocorreu um erro":                  scraper.ErrUnexpectedResponse,
		"object reference not set":         scraper.ErrUnexpectedResponse,
	}
	for pattern, err := range patterns {
		if strings.Contains(s, pattern) {
			return err
		}
	}
	return nil
}

func normalizeAccessKey(key string) string {
	var b strings.Builder
	for _, r := range key {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func toURLValues(m map[string]string) url.Values {
	v := url.Values{}
	for k, val := range m {
		v.Set(k, val)
	}
	return v
}
