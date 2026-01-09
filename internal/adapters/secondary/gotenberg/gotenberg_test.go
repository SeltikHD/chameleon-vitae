package gotenberg_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/gotenberg"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

func TestNew(t *testing.T) {
	t.Run("creates client with default config", func(t *testing.T) {
		cfg := gotenberg.Config{}
		client, err := gotenberg.New(cfg)
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("creates client with custom config", func(t *testing.T) {
		cfg := gotenberg.Config{
			URL:     "http://custom-gotenberg:3000",
			Timeout: 120 * 1e9,
		}
		client, err := gotenberg.New(cfg)
		require.NoError(t, err)
		require.NotNil(t, client)
	})
}

func TestDefaultConfig(t *testing.T) {
	cfg := gotenberg.DefaultConfig()
	assert.Equal(t, "http://localhost:3000", cfg.URL)
	assert.NotZero(t, cfg.Timeout)
}

func TestClose(t *testing.T) {
	cfg := gotenberg.Config{}
	client, err := gotenberg.New(cfg)
	require.NoError(t, err)
	err = client.Close()
	assert.NoError(t, err)
}

func TestGotenbergClientInterface(t *testing.T) {
	cfg := gotenberg.Config{}
	client, err := gotenberg.New(cfg)
	require.NoError(t, err)
	_ = client.GeneratePDF
	_ = client.GetTemplates
	_ = client.HealthCheck
	_ = client.Close
}

func TestGetTemplates(t *testing.T) {
	cfg := gotenberg.Config{}
	client, err := gotenberg.New(cfg)
	require.NoError(t, err)
	templates, err := client.GetTemplates(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, templates)
	for _, tmpl := range templates {
		assert.NotEmpty(t, tmpl.Name)
		assert.NotEmpty(t, tmpl.DisplayName)
	}
}

func TestGetTemplateByName(t *testing.T) {
	cfg := gotenberg.Config{}
	client, err := gotenberg.New(cfg)
	require.NoError(t, err)
	templates, err := client.GetTemplates(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, templates)

	t.Run("finds existing template", func(t *testing.T) {
		assert.Equal(t, "professional", templates[0].Name)
	})

	t.Run("has expected templates", func(t *testing.T) {
		names := make([]string, len(templates))
		for i, tmpl := range templates {
			names[i] = tmpl.Name
		}
		assert.Contains(t, names, "professional")
		assert.Contains(t, names, "minimal")
		assert.Contains(t, names, "technical")
	})
}

func TestGeneratePDFWithMockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("%PDF-1.4 mock pdf content"))
	}))
	defer server.Close()

	cfg := gotenberg.Config{URL: server.URL}
	client, err := gotenberg.New(cfg)
	require.NoError(t, err)
	ctx := context.Background()

	t.Run("generates PDF from HTML", func(t *testing.T) {
		result, err := client.GeneratePDF(ctx, ports.GeneratePDFRequest{
			HTML: "<html><body><h1>Test Resume</h1></body></html>",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		content, err := io.ReadAll(result.Content)
		require.NoError(t, err)
		assert.Contains(t, string(content), "%PDF")
		result.Content.Close()
	})

	t.Run("generates PDF with CSS", func(t *testing.T) {
		result, err := client.GeneratePDF(ctx, ports.GeneratePDFRequest{
			HTML: "<html><body><h1>Styled Resume</h1></body></html>",
			CSS:  "body { font-family: Arial; }",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		result.Content.Close()
	})

	t.Run("applies custom PDF options", func(t *testing.T) {
		result, err := client.GeneratePDF(ctx, ports.GeneratePDFRequest{
			HTML: "<html><body><h1>Custom Options</h1></body></html>",
			Options: ports.PDFOptions{
				PaperWidth:   8.5,
				PaperHeight:  11,
				MarginTop:    0.5,
				MarginBottom: 0.5,
				MarginLeft:   0.5,
				MarginRight:  0.5,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		result.Content.Close()
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("healthy server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				return
			}
			http.NotFound(w, r)
		}))
		defer server.Close()
		cfg := gotenberg.Config{URL: server.URL}
		client, err := gotenberg.New(cfg)
		require.NoError(t, err)
		err = client.HealthCheck(context.Background())
		assert.NoError(t, err)
	})

	t.Run("unhealthy server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer server.Close()
		cfg := gotenberg.Config{URL: server.URL}
		client, err := gotenberg.New(cfg)
		require.NoError(t, err)
		err = client.HealthCheck(context.Background())
		assert.Error(t, err)
	})

	t.Run("unreachable server", func(t *testing.T) {
		cfg := gotenberg.Config{URL: "http://localhost:99999"}
		client, err := gotenberg.New(cfg)
		require.NoError(t, err)
		err = client.HealthCheck(context.Background())
		assert.Error(t, err)
	})
}

func TestGeneratePDFErrors(t *testing.T) {
	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()
		cfg := gotenberg.Config{URL: server.URL}
		client, err := gotenberg.New(cfg)
		require.NoError(t, err)
		_, err = client.GeneratePDF(context.Background(), ports.GeneratePDFRequest{
			HTML: "<html><body>Test</body></html>",
		})
		require.Error(t, err)
	})
}
