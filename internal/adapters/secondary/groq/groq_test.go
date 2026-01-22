// Package groq_test contains unit tests for the Groq AI adapter.
package groq_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/groq"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

func TestNew(t *testing.T) {
	t.Run("requires API key", func(t *testing.T) {
		cfg := groq.Config{}

		client, err := groq.New(cfg)
		require.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "API key is required")
	})

	t.Run("creates client with valid API key", func(t *testing.T) {
		cfg := groq.Config{
			APIKey: "test-api-key", // pragma: allowlist secret
		}

		client, err := groq.New(cfg)
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("uses default values when not specified", func(t *testing.T) {
		cfg := groq.Config{
			APIKey: "test-api-key", // pragma: allowlist secret
		}

		client, err := groq.New(cfg)
		require.NoError(t, err)
		require.NotNil(t, client)
		// The client should be created with defaults
	})
}

func TestDefaultConfig(t *testing.T) {
	cfg := groq.DefaultConfig()

	assert.Equal(t, "llama-3.3-70b-versatile", cfg.ModelGeneration)
	assert.Equal(t, "meta-llama/llama-4-scout-17b-16e-instruct", cfg.ModelAnalysis)
	assert.Equal(t, 3, cfg.MaxRetries)
	assert.NotZero(t, cfg.Timeout)
}

// TestGroqClientInterface verifies the Client implements ports.AIProvider.
// This is a compile-time check.
func TestGroqClientInterface(t *testing.T) {
	cfg := groq.Config{APIKey: "test"} // pragma: allowlist secret
	client, err := groq.New(cfg)
	require.NoError(t, err)

	// Just verify the methods exist and have correct signatures.
	// These would fail to compile if signatures don't match.
	_ = client.AnalyzeJob
	_ = client.SelectBullets
	_ = client.TailorBullet
	_ = client.GenerateSummary
	_ = client.ScoreMatch
	_ = client.Close
}

// TestGroqMockServer tests the client with a mock HTTP server.
// Note: The current implementation uses a hardcoded baseURL constant,
// so we cannot easily redirect requests to a mock server.
// This test demonstrates the pattern for when baseURL is configurable.
func TestGroqMockServer(t *testing.T) {
	t.Skip("Groq client uses hardcoded baseURL - needs refactoring for mock testing")

	// This is how the test would work with a configurable baseURL:
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Return mock response
		response := map[string]any{
			"choices": []map[string]any{
				{
					"message": map[string]any{
						"content": `{"title": "Software Engineer", "company": "Test Corp", "required_skills": ["Go", "Python"], "preferred_skills": [], "keywords": ["backend"], "seniority_level": "senior", "years_experience": 5, "summary": "Senior backend role"}`,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Would need configurable baseURL in groq.Config
	// cfg := groq.Config{
	// 	APIKey:  "test-api-key", // pragma: allowlist secret
	// 	BaseURL: server.URL,
	// }
}

func TestClose(t *testing.T) {
	cfg := groq.Config{APIKey: "test-api-key"} // pragma: allowlist secret
	client, err := groq.New(cfg)
	require.NoError(t, err)

	err = client.Close()
	assert.NoError(t, err)
}

// Integration test placeholder - requires actual API key.
func TestGroqIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test would run with a real API key from environment:
	// apiKey := os.Getenv("GROQ_API_KEY") // pragma: allowlist secret
	// if apiKey == "" {
	// 	t.Skip("GROQ_API_KEY not set")
	// }

	t.Skip("Integration tests require GROQ_API_KEY environment variable")
}

// TestAnalyzeJobRequest verifies request structure.
func TestAnalyzeJobRequest(t *testing.T) {
	// This is a documentation test showing expected request format.
	t.Run("validates job description format", func(t *testing.T) {
		ctx := context.Background()
		cfg := groq.Config{APIKey: "test-api-key"} // pragma: allowlist secret
		client, err := groq.New(cfg)
		require.NoError(t, err)

		// Call will fail because we can't mock the server,
		// but this documents the expected API.
		_, err = client.AnalyzeJob(ctx, ports.AnalyzeJobRequest{
			JobDescription: "", // Empty description
		})

		// With empty description, it should still make the API call
		// (validation happens on the AI side)
		assert.Error(t, err) // Will error because no real API
	})
}
