// Package http provides an internal HTTP client for scraping operations.
package http

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type Client struct {
	http    *http.Client
	baseURL string
	headers map[string]string
}

type Option func(*Client)

func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.http.Timeout = d }
}

func WithHeader(key, value string) Option {
	return func(c *Client) { c.headers[key] = value }
}

func WithInsecureSkipVerify() Option {
	return func(c *Client) {
		c.http.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
}

func New(baseURL string, opts ...Option) (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	c := &Client{
		http:    &http.Client{Jar: jar, Timeout: 30 * time.Second},
		baseURL: strings.TrimSuffix(baseURL, "/"),
		headers: map[string]string{
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Accept-Language": "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7",
		},
	}

	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

type RequestConfig struct {
	Headers map[string]string
	Params  url.Values
	Referer string
}

type Response struct {
	*http.Response
	body []byte
}

func (r *Response) Body() ([]byte, error) {
	if r.body != nil {
		return r.body, nil
	}
	defer r.Response.Body.Close()
	body, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return nil, err
	}
	r.body = body
	return r.body, nil
}

func (c *Client) Get(ctx context.Context, path string, cfg *RequestConfig) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.buildURL(path, cfg), nil)
	if err != nil {
		return nil, err
	}
	c.applyHeaders(req, cfg)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	return &Response{Response: resp}, nil
}

func (c *Client) PostForm(ctx context.Context, path string, data url.Values, cfg *RequestConfig) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.buildURL(path, cfg), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.applyHeaders(req, cfg)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	return &Response{Response: resp}, nil
}

func (c *Client) GetImage(ctx context.Context, path string, cfg *RequestConfig) ([]byte, string, error) {
	resp, err := c.Get(ctx, path, cfg)
	if err != nil {
		return nil, "", err
	}
	body, err := resp.Body()
	if err != nil {
		return nil, "", err
	}
	return body, resp.Header.Get("Content-Type"), nil
}

func (c *Client) BaseURL() string { return c.baseURL }

func (c *Client) buildURL(path string, cfg *RequestConfig) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		if cfg != nil && len(cfg.Params) > 0 {
			return path + "?" + cfg.Params.Encode()
		}
		return path
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	fullURL := c.baseURL + path
	if cfg != nil && len(cfg.Params) > 0 {
		fullURL += "?" + cfg.Params.Encode()
	}
	return fullURL
}

func (c *Client) applyHeaders(req *http.Request, cfg *RequestConfig) {
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	if cfg != nil {
		for k, v := range cfg.Headers {
			req.Header.Set(k, v)
		}
		if cfg.Referer != "" {
			req.Header.Set("Referer", cfg.Referer)
		}
	}
}
