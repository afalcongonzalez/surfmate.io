package browser

import (
	"fmt"
	"sync"
	"time"

	"github.com/afalcongonzalez/surfmate.io/internal/config"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// Manager handles the browser instance and page operations
type Manager struct {
	browser *rod.Browser
	page    *rod.Page
	config  *config.Config
	mu      sync.Mutex
}

var (
	instance *Manager
	once     sync.Once
)

// GetManager returns the singleton browser manager
func GetManager(cfg *config.Config) *Manager {
	once.Do(func() {
		instance = &Manager{config: cfg}
	})
	return instance
}

// Launch starts the browser with the configured settings
func (m *Manager) Launch() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.browser != nil {
		return nil
	}

	l := launcher.New().
		Headless(m.config.Headless).
		Set("disable-gpu", "false").
		Set("no-sandbox", "").
		Set("disable-dev-shm-usage", "")

	// Use host browser if configured or detected
	browserPath := m.config.BrowserPath
	if browserPath == "" {
		browserPath = FindHostBrowser()
	}
	if browserPath != "" {
		l = l.Bin(browserPath)
	}

	url, err := l.Launch()
	if err != nil {
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	m.browser = rod.New().ControlURL(url)
	if err := m.browser.Connect(); err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}

	// Create initial page
	m.page, err = m.browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	// Set viewport
	err = m.page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:  m.config.ViewportWidth,
		Height: m.config.ViewportHeight,
	})
	if err != nil {
		return fmt.Errorf("failed to set viewport: %w", err)
	}

	return nil
}

// Close shuts down the browser
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.browser != nil {
		return m.browser.Close()
	}
	return nil
}

// Page returns the current page
func (m *Manager) Page() *rod.Page {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.page
}

// Navigate goes to the specified URL
func (m *Manager) Navigate(url string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.page.Timeout(m.config.BrowserTimeout).Navigate(url)
}

// WaitLoad waits for the page to finish loading
func (m *Manager) WaitLoad() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.page.Timeout(m.config.BrowserTimeout).WaitLoad()
}

// GetTitle returns the page title
func (m *Manager) GetTitle() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	info, err := m.page.Info()
	if err != nil {
		return "", err
	}
	return info.Title, nil
}

// GetURL returns the current page URL
func (m *Manager) GetURL() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	info, err := m.page.Info()
	if err != nil {
		return "", err
	}
	return info.URL, nil
}

// Click clicks on an element matching the selector
func (m *Manager) Click(selector string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	el, err := m.page.Timeout(m.config.BrowserTimeout).Element(selector)
	if err != nil {
		return fmt.Errorf("element not found: %s", selector)
	}
	return el.Click(proto.InputMouseButtonLeft, 1)
}

// Type types text into an element matching the selector
func (m *Manager) Type(selector, text string, submit bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	el, err := m.page.Timeout(m.config.BrowserTimeout).Element(selector)
	if err != nil {
		return fmt.Errorf("element not found: %s", selector)
	}

	if err := el.Input(text); err != nil {
		return err
	}

	if submit {
		return m.page.Keyboard.Press(13) // Enter key
	}
	return nil
}

// Scroll scrolls the page in the specified direction
func (m *Manager) Scroll(direction string, amount int, selector string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if selector != "" {
		el, err := m.page.Timeout(m.config.BrowserTimeout).Element(selector)
		if err != nil {
			return fmt.Errorf("element not found: %s", selector)
		}
		return el.ScrollIntoView()
	}

	var deltaX, deltaY float64
	switch direction {
	case "up":
		deltaY = float64(-amount)
	case "down":
		deltaY = float64(amount)
	case "left":
		deltaX = float64(-amount)
	case "right":
		deltaX = float64(amount)
	default:
		deltaY = float64(amount)
	}

	return m.page.Mouse.Scroll(deltaX, deltaY, 1)
}

// GetContent returns the page content
func (m *Manager) GetContent(includeHTML bool) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if includeHTML {
		return m.page.HTML()
	}
	el, err := m.page.Element("body")
	if err != nil {
		return "", err
	}
	return el.Text()
}

// Screenshot captures the page as a base64 encoded image
func (m *Manager) Screenshot(fullPage bool, selector string, quality int) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if selector != "" {
		el, err := m.page.Timeout(m.config.BrowserTimeout).Element(selector)
		if err != nil {
			return nil, fmt.Errorf("element not found: %s", selector)
		}
		return el.Screenshot(proto.PageCaptureScreenshotFormatPng, quality)
	}

	if fullPage {
		return m.page.Screenshot(true, &proto.PageCaptureScreenshot{
			Format:  proto.PageCaptureScreenshotFormatPng,
			Quality: &quality,
		})
	}

	return m.page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatPng,
		Quality: &quality,
	})
}

// ExtractText extracts text from elements matching the selector
func (m *Manager) ExtractText(selector string, multiple bool) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if multiple {
		elements, err := m.page.Timeout(m.config.BrowserTimeout).Elements(selector)
		if err != nil {
			return nil, fmt.Errorf("elements not found: %s", selector)
		}
		var texts []string
		for _, el := range elements {
			text, err := el.Text()
			if err == nil {
				texts = append(texts, text)
			}
		}
		return texts, nil
	}

	el, err := m.page.Timeout(m.config.BrowserTimeout).Element(selector)
	if err != nil {
		return nil, fmt.Errorf("element not found: %s", selector)
	}
	text, err := el.Text()
	if err != nil {
		return nil, err
	}
	return []string{text}, nil
}

// WaitForUser waits for user intervention (e.g., captcha solving)
func (m *Manager) WaitForUser(timeout time.Duration) error {
	m.mu.Lock()
	page := m.page
	m.mu.Unlock()

	return WaitForCaptchaResolution(page, timeout)
}

// HasCaptcha checks if a captcha is present on the page
func (m *Manager) HasCaptcha() bool {
	m.mu.Lock()
	page := m.page
	m.mu.Unlock()

	return DetectCaptcha(page)
}
