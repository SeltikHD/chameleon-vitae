// Package gotenberg provides a PDF generation adapter using Gotenberg.
package gotenberg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

const (
	chromiumEndpoint = "/forms/chromium/convert/html"
	healthEndpoint   = "/health"
)

// Config holds Gotenberg configuration.
type Config struct {
	// URL is the Gotenberg service URL.
	URL string

	// Timeout is the HTTP request timeout.
	Timeout time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		URL:     "http://localhost:3000",
		Timeout: 60 * time.Second,
	}
}

// Client implements ports.PDFEngine using Gotenberg.
type Client struct {
	config     Config
	httpClient *http.Client
	templates  []ports.PDFTemplate
}

// New creates a new Gotenberg client.
func New(cfg Config) (*Client, error) {
	if cfg.URL == "" {
		cfg.URL = DefaultConfig().URL
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultConfig().Timeout
	}

	client := &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		templates: defaultTemplates(),
	}

	return client, nil
}

// GeneratePDF generates a PDF from HTML content.
func (c *Client) GeneratePDF(ctx context.Context, req ports.GeneratePDFRequest) (*ports.PDFResult, error) {
	// Prepare multipart form.
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add HTML file.
	htmlContent := req.HTML
	if req.CSS != "" {
		// Inject CSS into HTML if provided.
		htmlContent = injectCSS(htmlContent, req.CSS)
	}

	htmlPart, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, fmt.Errorf("gotenberg: failed to create form file: %w", err)
	}
	if _, err := htmlPart.Write([]byte(htmlContent)); err != nil {
		return nil, fmt.Errorf("gotenberg: failed to write HTML: %w", err)
	}

	// Apply PDF options.
	opts := req.Options
	if opts.PaperWidth == 0 {
		opts = ports.DefaultPDFOptions()
	}

	// Add form fields for PDF options.
	formFields := map[string]string{
		"paperWidth":        strconv.FormatFloat(opts.PaperWidth, 'f', 2, 64),
		"paperHeight":       strconv.FormatFloat(opts.PaperHeight, 'f', 2, 64),
		"marginTop":         strconv.FormatFloat(opts.MarginTop, 'f', 2, 64),
		"marginBottom":      strconv.FormatFloat(opts.MarginBottom, 'f', 2, 64),
		"marginLeft":        strconv.FormatFloat(opts.MarginLeft, 'f', 2, 64),
		"marginRight":       strconv.FormatFloat(opts.MarginRight, 'f', 2, 64),
		"scale":             strconv.FormatFloat(opts.Scale, 'f', 2, 64),
		"printBackground":   "true",
		"preferCssPageSize": "false",
	}

	for key, value := range formFields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, fmt.Errorf("gotenberg: failed to write field %s: %w", key, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("gotenberg: failed to close writer: %w", err)
	}

	// Create request.
	url := c.config.URL + chromiumEndpoint
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("gotenberg: failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute request.
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("gotenberg: request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gotenberg: conversion failed (status %d): %s", resp.StatusCode, string(body))
	}

	// Generate filename.
	filename := "resume.pdf"
	if req.TemplateName != "" {
		filename = fmt.Sprintf("resume_%s.pdf", req.TemplateName)
	}

	return &ports.PDFResult{
		Content:  resp.Body, // Caller is responsible for closing.
		Size:     resp.ContentLength,
		Filename: filename,
	}, nil
}

// GetTemplates returns available resume templates.
func (c *Client) GetTemplates(ctx context.Context) ([]ports.PDFTemplate, error) {
	return c.templates, nil
}

// HealthCheck checks if the PDF engine is available.
func (c *Client) HealthCheck(ctx context.Context) error {
	url := c.config.URL + healthEndpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("gotenberg: failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gotenberg: health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gotenberg: unhealthy (status %d)", resp.StatusCode)
	}

	return nil
}

// Close releases any resources held by the PDF engine.
func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// injectCSS injects CSS into HTML head.
func injectCSS(html, css string) string {
	styleTag := fmt.Sprintf("<style>%s</style>", css)

	// Try to inject before </head>.
	if idx := bytes.Index([]byte(html), []byte("</head>")); idx != -1 {
		return html[:idx] + styleTag + html[idx:]
	}

	// Fallback: prepend to HTML.
	return styleTag + html
}

// defaultTemplates returns the built-in resume templates.
func defaultTemplates() []ports.PDFTemplate {
	return []ports.PDFTemplate{
		{
			Name:        "jake",
			DisplayName: "Jake's Resume",
			Description: "Industry gold standard for developer resumes. Single-page, dense, ATS-friendly format.",
			PreviewURL:  "/templates/jake/preview.png",
		},
		{
			Name:        "professional",
			DisplayName: "Professional",
			Description: "Clean, modern design suitable for most industries. Emphasizes readability and structure.",
			PreviewURL:  "/templates/professional/preview.png",
		},
		{
			Name:        "minimal",
			DisplayName: "Minimal",
			Description: "Simple, elegant layout with plenty of white space. Perfect for creative roles.",
			PreviewURL:  "/templates/minimal/preview.png",
		},
		{
			Name:        "technical",
			DisplayName: "Technical",
			Description: "Optimized for engineering and tech roles. Highlights skills and projects prominently.",
			PreviewURL:  "/templates/technical/preview.png",
		},
		{
			Name:        "executive",
			DisplayName: "Executive",
			Description: "Sophisticated design for senior leadership positions. Emphasizes achievements and impact.",
			PreviewURL:  "/templates/executive/preview.png",
		},
		{
			Name:        "academic",
			DisplayName: "Academic",
			Description: "Traditional CV format for research and academic positions. Supports publications and grants.",
			PreviewURL:  "/templates/academic/preview.png",
		},
	}
}
