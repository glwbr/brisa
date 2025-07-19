package juice

// captcha.go - Captcha handling strategies
// package juice
//
// import (
// 	"context"
// 	"fmt"
// 	"image"
// 	"io"
// 	"net/url"
// 	"strconv"
// 	"time"
// )
//
// // CaptchaMode defines how captchas are handled
// type CaptchaMode int
//
// const (
// 	// CaptchaModeManual returns the captcha for manual solving
// 	CaptchaModeManual CaptchaMode = iota
//
// 	// CaptchaModeOCR attempts to solve captcha with built-in OCR
// 	CaptchaModeOCR
//
// 	// CaptchaModeService uses a third-party service who knows?
// 	CaptchaModeService
// )
//
// // CaptchaResult represents a captcha challenge
// type CaptchaResult struct {
// 	Image       image.Image // For manual mode
// 	ImageBytes  []byte      // Raw image bytes
// 	Solution    string      // For automated modes
// 	NeedsManual bool        // Indicates if user input is needed
// }
//
// // CaptchaSolver interface for different solving strategies
// type CaptchaSolver interface {
// 	Solve(captchaImg io.Reader) (*CaptchaResult, error)
// }
//
// // GetCAPTCHA retrieves/generates a new captcha image
// func (c *Client) GetCAPTCHA(ctx context.Context) (string, error) {
// 	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
//
// 	res, err := c.Get(ctx, CaptchaPath, &RequestOptions{
// 		QueryParams: url.Values{
// 			"t": {strconv.FormatInt(timestamp, 10)},
// 		},
// 	})
// 	if err != nil {
// 		// TODO: specify new errors
// 		return "", fmt.Errorf("[juicy]: error generating the captcha: %w", err)
// 	}
// 	defer res.Body.Close()
//
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("[juicy]: error reading body: %v", err)
// 	}
//
// 	return string(body), nil
// }
//
// // manualCaptchaSolver implements CaptchaSolver for manual solving
// type manualCaptchaSolver struct{}
//
// func (s *manualCaptchaSolver) Solve(captchaImg io.Reader) (*CaptchaResult, error) {
// 	imgBytes, err := io.ReadAll(captchaImg)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &CaptchaResult{
// 		ImageBytes:  imgBytes,
// 		NeedsManual: true,
// 	}, nil
// }
