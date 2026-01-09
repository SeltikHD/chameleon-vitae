// Package jina provides a job parsing adapter using the Jina Reader API.
package jina

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

const (
	readerBaseURL = "https://r.jina.ai"
)

// Config holds Jina API configuration.
type Config struct {
	// APIKey is the Jina API key.
	APIKey string

	// Timeout is the HTTP request timeout.
	Timeout time.Duration

	// MaxRetries is the maximum number of retries on errors.
	MaxRetries int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}
}

// Client implements ports.JobParser using the Jina Reader API.
type Client struct {
	config     Config
	httpClient *http.Client
}

// New creates a new Jina API client.
func New(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("jina: API key is required")
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultConfig().Timeout
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = DefaultConfig().MaxRetries
	}

	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

// ParseJobURL fetches and parses a job description from a URL.
func (c *Client) ParseJobURL(ctx context.Context, jobURL string) (*ports.ParsedJob, error) {
	// Validate URL.
	parsedURL, err := url.Parse(jobURL)
	if err != nil {
		return nil, fmt.Errorf("jina: invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("jina: URL must use http or https scheme")
	}

	// Build Jina Reader URL.
	readerURL := fmt.Sprintf("%s/%s", readerBaseURL, jobURL)

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff.
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		result, err := c.doRequest(ctx, readerURL)
		if err != nil {
			lastErr = err
			continue
		}

		return &ports.ParsedJob{
			URL:           jobURL,
			Title:         result.Title,
			Content:       result.Content,
			Description:   extractJobDescription(result.Content),
			PublishedDate: result.PublishedDate,
			Metadata:      result.Metadata,
		}, nil
	}

	return nil, fmt.Errorf("jina: max retries exceeded: %w", lastErr)
}

// HealthCheck checks if the parser is available.
func (c *Client) HealthCheck(ctx context.Context) error {
	// Jina doesn't have a dedicated health endpoint.
	// We'll try to parse a simple URL as a health check.
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, readerBaseURL, nil)
	if err != nil {
		return fmt.Errorf("jina: failed to create health check request: %w", err)
	}

	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("jina: health check failed: %w", err)
	}
	defer resp.Body.Close()

	// Accept any 2xx or 4xx response as "healthy" (service is responding).
	if resp.StatusCode >= 500 {
		return fmt.Errorf("jina: unhealthy (status %d)", resp.StatusCode)
	}

	return nil
}

// Close releases any resources held by the job parser.
func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// readerResult represents the response from Jina Reader API.
type readerResult struct {
	Title         string
	Content       string
	PublishedDate string
	Metadata      map[string]string
}

// doRequest performs the HTTP request to Jina Reader.
func (c *Client) doRequest(ctx context.Context, url string) (*readerResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers.
	req.Header.Set("Accept", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Try to parse as JSON first.
	var jsonResp struct {
		Code   int `json:"code"`
		Status int `json:"status"`
		Data   struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			Content     string `json:"content"`
			PublishTime string `json:"publishTime"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &jsonResp); err == nil && jsonResp.Data.Content != "" {
		metadata := make(map[string]string)
		if jsonResp.Data.Description != "" {
			metadata["description"] = jsonResp.Data.Description
		}

		return &readerResult{
			Title:         jsonResp.Data.Title,
			Content:       jsonResp.Data.Content,
			PublishedDate: jsonResp.Data.PublishTime,
			Metadata:      metadata,
		}, nil
	}

	// Fallback: treat response as plain markdown text.
	content := string(body)
	title := extractTitle(content)

	return &readerResult{
		Title:    title,
		Content:  content,
		Metadata: make(map[string]string),
	}, nil
}

// extractTitle attempts to extract a title from markdown content.
func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}

	// Fallback: return first non-empty line.
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// Truncate if too long.
			if len(line) > 100 {
				return line[:100] + "..."
			}
			return line
		}
	}

	return "Job Description"
}

// extractJobDescription attempts to extract the main job description from content.
func extractJobDescription(content string) string {
	// For now, return the content as-is.
	// Future: implement smarter extraction logic to identify the job description section.

	// Remove excessive whitespace.
	lines := strings.Split(content, "\n")
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" || (len(cleanLines) > 0 && cleanLines[len(cleanLines)-1] != "") {
			cleanLines = append(cleanLines, trimmed)
		}
	}

	return strings.Join(cleanLines, "\n")
}
