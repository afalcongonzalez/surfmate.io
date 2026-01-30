package browser

import (
	"time"

	"github.com/go-rod/rod"
)

// CaptchaSelectors contains CSS selectors for common captcha implementations
var CaptchaSelectors = []string{
	".g-recaptcha",
	".h-captcha",
	".cf-turnstile",
	"iframe[src*='captcha']",
	"iframe[src*='challenge']",
	"iframe[src*='recaptcha']",
	"iframe[src*='hcaptcha']",
	"[data-sitekey]",
	"#px-captcha",
	".captcha",
	"[class*='captcha']",
	"[id*='captcha']",
}

// DetectCaptcha checks if a captcha is present on the page
func DetectCaptcha(page *rod.Page) bool {
	for _, selector := range CaptchaSelectors {
		el, err := page.Timeout(500 * time.Millisecond).Element(selector)
		if err == nil && el != nil {
			visible, _ := el.Visible()
			if visible {
				return true
			}
		}
	}
	return false
}

// WaitForCaptchaResolution polls until the captcha element disappears or timeout
func WaitForCaptchaResolution(page *rod.Page, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	pollInterval := 500 * time.Millisecond

	for time.Now().Before(deadline) {
		if !DetectCaptcha(page) {
			return nil
		}
		time.Sleep(pollInterval)
	}

	return nil // Return nil even on timeout - user may have navigated away
}
