// Package scraper provides HTTP client and scraping utilities for NFC-e portals.
package scraper

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

// Client is an HTTP client with cookie jar support for stateful scraping.
type Client struct {
	http    *http.Client
	baseURL string
	headers map[string]string
}

// ClientOption configures the Client.
type ClientOption func(*Client)

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		c.http.Timeout = d
	}
}

// WithHeader sets a default header for all requests.
func WithHeader(key, value string) ClientOption {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// WithCookieJar sets a custom cookie jar.
func WithCookieJar(jar http.CookieJar) ClientOption {
	return func(c *Client) {
		c.http.Jar = jar
	}
}

// WithInsecureSkipVerify disables TLS certificate verification.
// Use this only for portals with self-signed or problematic certificates.
func WithInsecureSkipVerify() ClientOption {
	return func(c *Client) {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		c.http.Transport = transport
	}
}

// NewClient creates a new scraping HTTP client with cookie support.
func NewClient(baseURL string, opts ...ClientOption) (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}

	c := &Client{
		http: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
		baseURL: strings.TrimSuffix(baseURL, "/"),
		headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// RequestConfig holds optional request configuration.
type RequestConfig struct {
	// Headers are additional headers to send with this request.
	Headers map[string]string
	// Params are URL query parameters.
	Params url.Values
	// Referer sets the Referer header.
	Referer string
}

// Response wraps an HTTP response with convenience methods.
type Response struct {
	*http.Response
	body []byte
}

// Body returns the response body as bytes, reading it lazily.
func (r *Response) Body() ([]byte, error) {
	if r.body != nil {
		return r.body, nil
	}

	defer r.Response.Body.Close()
	body, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	r.body = body
	return r.body, nil
}

// Get performs an HTTP GET request.
func (c *Client) Get(ctx context.Context, path string, cfg *RequestConfig) (*Response, error) {
	fullURL := c.buildURL(path, cfg)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.applyHeaders(req, cfg)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	return &Response{Response: resp}, nil
}

// PostForm performs an HTTP POST request with form data.
func (c *Client) PostForm(ctx context.Context, path string, data url.Values, cfg *RequestConfig) (*Response, error) {
	fullURL := c.buildURL(path, cfg)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.applyHeaders(req, cfg)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	return &Response{Response: resp}, nil
}

// GetImage performs an HTTP GET request expecting an image response.
func (c *Client) GetImage(ctx context.Context, path string, cfg *RequestConfig) ([]byte, string, error) {
	resp, err := c.Get(ctx, path, cfg)
	if err != nil {
		return nil, "", err
	}

	body, err := resp.Body()
	if err != nil {
		return nil, "", err
	}

	contentType := resp.Header.Get("Content-Type")
	return body, contentType, nil
}

// buildURL constructs the full URL from base URL, path, and query params.
func (c *Client) buildURL(path string, cfg *RequestConfig) string {
	// If path is already a full URL, use it directly
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		if cfg != nil && len(cfg.Params) > 0 {
			return path + "?" + cfg.Params.Encode()
		}
		return path
	}

	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	fullURL := c.baseURL + path
	if cfg != nil && len(cfg.Params) > 0 {
		fullURL += "?" + cfg.Params.Encode()
	}

	return fullURL
}

// applyHeaders applies default and request-specific headers.
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

// Cookies returns the cookies for the given URL.
func (c *Client) Cookies(u *url.URL) []*http.Cookie {
	return c.http.Jar.Cookies(u)
}

// SetCookies sets cookies for the given URL.
func (c *Client) SetCookies(u *url.URL, cookies []*http.Cookie) {
	c.http.Jar.SetCookies(u, cookies)
}

// BaseURL returns the client's base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}
