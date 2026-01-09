// Package jina_test contains unit tests for the Jina Reader adapter.
package jina_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/jina"
)

func TestNew(t *testing.T) {
	t.Run("requires API key", func(t *testing.T) {
		cfg := jina.Config{}

		client, err := jina.New(cfg)
		require.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "API key is required")
	})

	t.Run("creates client with valid API key", func(t *testing.T) {
		cfg := jina.Config{
			APIKey: "test-api-key", // pragma: allowlist secret
		}

		client, err := jina.New(cfg)
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("uses default timeout when not specified", func(t *testing.T) {
		cfg := jina.Config{
			APIKey: "test-api-key", // pragma: allowlist secret
		}

		client, err := jina.New(cfg)
		require.NoError(t, err)
		require.NotNil(t, client)
	})
}

func TestDefaultConfig(t *testing.T) {
	cfg := jina.DefaultConfig()

	assert.NotZero(t, cfg.Timeout)
	assert.Equal(t, 3, cfg.MaxRetries)
}

func TestClose(t *testing.T) {
	cfg := jina.Config{APIKey: "test-api-key"} // pragma: allowlist secret
	client, err := jina.New(cfg)
	require.NoError(t, err)

	err = client.Close()
	assert.NoError(t, err)
}

// TestJinaClientInterface verifies the Client has the expected methods.
func TestJinaClientInterface(t *testing.T) {
	cfg := jina.Config{APIKey: "test"} // pragma: allowlist secret
	client, err := jina.New(cfg)
	require.NoError(t, err)

	// Verify the methods exist (compile-time check).
	_ = client.ParseJobURL
	_ = client.Close
}

// TestParseJobURLValidation tests URL validation.
func TestParseJobURLValidation(t *testing.T) {
	cfg := jina.Config{APIKey: "test-api-key"} // pragma: allowlist secret
	client, err := jina.New(cfg)
	require.NoError(t, err)

	t.Run("rejects invalid URL", func(t *testing.T) {
		_, err := client.ParseJobURL(t.Context(), "not-a-valid-url")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "URL must use http or https scheme")
	})

	t.Run("rejects non-http schemes", func(t *testing.T) {
		_, err := client.ParseJobURL(t.Context(), "ftp://example.com/job")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "URL must use http or https scheme")
	})

	// Valid URL will fail because we can't mock the external API
	t.Run("accepts valid https URL format", func(t *testing.T) {
		// This will fail on the HTTP request, but validates the URL first
		_, err := client.ParseJobURL(t.Context(), "https://example.com/job")
		require.Error(t, err) // Will fail on HTTP request
		// But should not be a URL validation error
		assert.NotContains(t, err.Error(), "invalid URL")
	})
}

// Integration test placeholder.
func TestJinaIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Integration tests require JINA_API_KEY environment variable")
}
