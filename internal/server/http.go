package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/afalcongonzalez/surfmate.io/internal/browser"
)

// HTTPServer provides REST API endpoints for browser automation
type HTTPServer struct {
	mgr  *browser.Manager
	port int
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(mgr *browser.Manager, port int) *HTTPServer {
	return &HTTPServer{mgr: mgr, port: port}
}

// response helpers
func jsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Start runs the HTTP server
func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()

	// OpenAPI spec
	mux.HandleFunc("GET /openapi.yaml", s.handleOpenAPI)
	mux.HandleFunc("GET /", s.handleRoot)

	// Browser tools
	mux.HandleFunc("POST /open_browser", s.handleOpenBrowser)
	mux.HandleFunc("POST /navigate", s.handleNavigate)
	mux.HandleFunc("POST /click", s.handleClick)
	mux.HandleFunc("POST /type", s.handleType)
	mux.HandleFunc("POST /scroll", s.handleScroll)
	mux.HandleFunc("GET /content", s.handleGetContent)
	mux.HandleFunc("GET /screenshot", s.handleScreenshot)
	mux.HandleFunc("POST /extract", s.handleExtract)
	mux.HandleFunc("POST /wait", s.handleWait)

	// CORS middleware
	handler := corsMiddleware(mux)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Starting HTTP server on http://localhost%s\n", addr)
	fmt.Printf("OpenAPI spec: http://localhost%s/openapi.yaml\n", addr)
	return http.ListenAndServe(addr, handler)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *HTTPServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, map[string]string{
		"name":    "surfmate.io",
		"version": "1.0.0",
		"docs":    "/openapi.yaml",
	})
}

func (s *HTTPServer) handleOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml")
	w.Write([]byte(openAPISpec))
}

func (s *HTTPServer) handleOpenBrowser(w http.ResponseWriter, r *http.Request) {
	if s.mgr.IsLaunched() {
		jsonResponse(w, map[string]string{"status": "already_open", "message": "Browser is already open."})
		return
	}

	if err := s.mgr.Launch(); err != nil {
		errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to launch browser: %v", err))
		return
	}

	jsonResponse(w, map[string]string{"status": "launched", "message": "Browser launched successfully."})
}

// NavigateRequest is the request body for /navigate
type NavigateRequest struct {
	URL string `json:"url"`
}

func (s *HTTPServer) handleNavigate(w http.ResponseWriter, r *http.Request) {
	var req NavigateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.URL == "" {
		errorResponse(w, http.StatusBadRequest, "url is required")
		return
	}

	if err := s.mgr.Navigate(req.URL); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := s.mgr.WaitLoad(); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	title, _ := s.mgr.GetTitle()
	currentURL, _ := s.mgr.GetURL()
	hasCaptcha := s.mgr.HasCaptcha()

	jsonResponse(w, map[string]any{
		"url":            currentURL,
		"title":          title,
		"captcha_found":  hasCaptcha,
	})
}

// ClickRequest is the request body for /click
type ClickRequest struct {
	Selector string `json:"selector"`
}

func (s *HTTPServer) handleClick(w http.ResponseWriter, r *http.Request) {
	var req ClickRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Selector == "" {
		errorResponse(w, http.StatusBadRequest, "selector is required")
		return
	}

	if err := s.mgr.Click(req.Selector); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, map[string]string{"status": "clicked", "selector": req.Selector})
}

// TypeRequest is the request body for /type
type TypeRequest struct {
	Selector string `json:"selector"`
	Text     string `json:"text"`
	Submit   bool   `json:"submit"`
}

func (s *HTTPServer) handleType(w http.ResponseWriter, r *http.Request) {
	var req TypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Selector == "" {
		errorResponse(w, http.StatusBadRequest, "selector is required")
		return
	}

	if err := s.mgr.Type(req.Selector, req.Text, req.Submit); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, map[string]any{
		"status":    "typed",
		"selector":  req.Selector,
		"submitted": req.Submit,
	})
}

// ScrollRequest is the request body for /scroll
type ScrollRequest struct {
	Direction string `json:"direction"`
	Amount    int    `json:"amount"`
	Selector  string `json:"selector"`
}

func (s *HTTPServer) handleScroll(w http.ResponseWriter, r *http.Request) {
	var req ScrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Direction == "" {
		req.Direction = "down"
	}
	if req.Amount == 0 {
		req.Amount = 300
	}

	if err := s.mgr.Scroll(req.Direction, req.Amount, req.Selector); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, map[string]any{
		"status":    "scrolled",
		"direction": req.Direction,
		"amount":    req.Amount,
	})
}

func (s *HTTPServer) handleGetContent(w http.ResponseWriter, r *http.Request) {
	includeHTML := r.URL.Query().Get("include_html") == "true"

	content, err := s.mgr.GetContent(includeHTML)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, map[string]string{"content": content})
}

func (s *HTTPServer) handleScreenshot(w http.ResponseWriter, r *http.Request) {
	fullPage := r.URL.Query().Get("full_page") == "true"
	selector := r.URL.Query().Get("selector")
	quality := 80

	data, err := s.mgr.Screenshot(fullPage, selector, quality)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	jsonResponse(w, map[string]string{
		"image":    encoded,
		"mimeType": "image/png",
	})
}

// ExtractRequest is the request body for /extract
type ExtractRequest struct {
	Selector string `json:"selector"`
	Multiple bool   `json:"multiple"`
}

func (s *HTTPServer) handleExtract(w http.ResponseWriter, r *http.Request) {
	var req ExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Selector == "" {
		errorResponse(w, http.StatusBadRequest, "selector is required")
		return
	}

	texts, err := s.mgr.ExtractText(req.Selector, req.Multiple)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, map[string]any{"texts": texts})
}

// WaitRequest is the request body for /wait
type WaitRequest struct {
	Reason  string `json:"reason"`
	Timeout int    `json:"timeout"`
}

func (s *HTTPServer) handleWait(w http.ResponseWriter, r *http.Request) {
	var req WaitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	timeout := 300 * time.Second
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
	}

	if err := s.mgr.WaitForUser(timeout); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	hasCaptcha := s.mgr.HasCaptcha()
	jsonResponse(w, map[string]any{
		"status":          "completed",
		"captcha_present": hasCaptcha,
	})
}

const openAPISpec = `openapi: 3.1.0
info:
  title: Surfmate.io Browser Automation
  description: Control a browser for web automation. The browser window is visible so users can solve captchas and complete logins.
  version: 1.0.0
servers:
  - url: http://localhost:8080
paths:
  /open_browser:
    post:
      operationId: openBrowser
      summary: Open the browser
      description: Launches the browser window. Must be called before using any other browser tools.
      responses:
        '200':
          description: Browser opened
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    enum: [launched, already_open]
                  message:
                    type: string

  /navigate:
    post:
      operationId: navigate
      summary: Navigate to a URL
      description: Opens a URL in the browser and returns the page title and whether a captcha was detected.
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [url]
              properties:
                url:
                  type: string
                  description: The URL to navigate to
      responses:
        '200':
          description: Navigation successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  url:
                    type: string
                  title:
                    type: string
                  captcha_found:
                    type: boolean

  /click:
    post:
      operationId: click
      summary: Click an element
      description: Clicks on an element matching the CSS selector.
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [selector]
              properties:
                selector:
                  type: string
                  description: CSS selector for the element to click
      responses:
        '200':
          description: Click successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                  selector:
                    type: string

  /type:
    post:
      operationId: typeText
      summary: Type text into an input
      description: Types text into an input element. Optionally presses Enter to submit.
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [selector]
              properties:
                selector:
                  type: string
                  description: CSS selector for the input element
                text:
                  type: string
                  description: Text to type
                submit:
                  type: boolean
                  description: Press Enter after typing
      responses:
        '200':
          description: Typing successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                  submitted:
                    type: boolean

  /scroll:
    post:
      operationId: scroll
      summary: Scroll the page
      description: Scrolls the page in a direction or scrolls an element into view.
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                direction:
                  type: string
                  enum: [up, down, left, right]
                  default: down
                amount:
                  type: integer
                  default: 300
                  description: Pixels to scroll
                selector:
                  type: string
                  description: If provided, scroll this element into view instead
      responses:
        '200':
          description: Scroll successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string

  /content:
    get:
      operationId: getContent
      summary: Get page content
      description: Returns the text content or HTML of the current page.
      parameters:
        - name: include_html
          in: query
          schema:
            type: boolean
            default: false
          description: Return full HTML instead of text
      responses:
        '200':
          description: Content retrieved
          content:
            application/json:
              schema:
                type: object
                properties:
                  content:
                    type: string

  /screenshot:
    get:
      operationId: screenshot
      summary: Take a screenshot
      description: Captures a screenshot of the page as a base64 PNG.
      parameters:
        - name: full_page
          in: query
          schema:
            type: boolean
            default: false
          description: Capture the full scrollable page
        - name: selector
          in: query
          schema:
            type: string
          description: Capture only this element
      responses:
        '200':
          description: Screenshot captured
          content:
            application/json:
              schema:
                type: object
                properties:
                  image:
                    type: string
                    description: Base64 encoded PNG
                  mimeType:
                    type: string

  /extract:
    post:
      operationId: extractText
      summary: Extract text from elements
      description: Extracts text content from elements matching a selector.
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [selector]
              properties:
                selector:
                  type: string
                  description: CSS selector for elements
                multiple:
                  type: boolean
                  default: false
                  description: Extract from all matching elements
      responses:
        '200':
          description: Text extracted
          content:
            application/json:
              schema:
                type: object
                properties:
                  texts:
                    type: array
                    items:
                      type: string

  /wait:
    post:
      operationId: waitForUser
      summary: Wait for user action
      description: Pauses and waits for user to complete an action like solving a captcha or logging in. Polls until captcha disappears.
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                reason:
                  type: string
                  description: Reason for waiting (shown to user)
                timeout:
                  type: integer
                  default: 300
                  description: Max seconds to wait
      responses:
        '200':
          description: Wait completed
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                  captcha_present:
                    type: boolean
`
