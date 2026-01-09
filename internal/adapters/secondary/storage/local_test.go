// Package storage_test contains unit tests for the storage adapters.
package storage_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/adapters/secondary/storage"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

func TestNewLocalStorage(t *testing.T) {
	t.Run("creates storage with default config", func(t *testing.T) {
		s, err := storage.NewLocalStorage(storage.LocalConfig{})
		require.NoError(t, err)
		require.NotNil(t, s)

		defer func() {
			_ = os.RemoveAll("./storage")
		}()
	})

	t.Run("creates storage with custom config", func(t *testing.T) {
		tempDir := t.TempDir()

		cfg := storage.LocalConfig{
			BasePath: tempDir,
			BaseURL:  "http://example.com/files",
		}

		s, err := storage.NewLocalStorage(cfg)
		require.NoError(t, err)
		require.NotNil(t, s)
	})

	t.Run("creates necessary directories", func(t *testing.T) {
		tempDir := t.TempDir()
		subDir := filepath.Join(tempDir, "nested", "storage")

		cfg := storage.LocalConfig{
			BasePath: subDir,
			BaseURL:  "http://example.com/files",
		}

		s, err := storage.NewLocalStorage(cfg)
		require.NoError(t, err)
		require.NotNil(t, s)

		info, err := os.Stat(subDir)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})
}

func TestLocalStorageUpload(t *testing.T) {
	tempDir := t.TempDir()
	baseURL := "http://example.com/files"

	s, err := storage.NewLocalStorage(storage.LocalConfig{
		BasePath: tempDir,
		BaseURL:  baseURL,
	})
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("uploads file with provided key", func(t *testing.T) {
		content := []byte("Hello, World!")
		key := "test-file.txt"

		result, err := s.Upload(ctx, ports.UploadRequest{
			Key:         key,
			Content:     bytes.NewReader(content),
			ContentType: "text/plain",
		})
		require.NoError(t, err)

		assert.Equal(t, key, result.Key)
		assert.Equal(t, baseURL+"/"+key, result.URL)
		assert.Equal(t, int64(len(content)), result.Size)

		data, err := os.ReadFile(filepath.Join(tempDir, key))
		require.NoError(t, err)
		assert.Equal(t, content, data)
	})

	t.Run("generates UUID key if not provided", func(t *testing.T) {
		content := []byte("Auto-keyed content")

		result, err := s.Upload(ctx, ports.UploadRequest{
			Content:     bytes.NewReader(content),
			ContentType: "text/plain",
		})
		require.NoError(t, err)

		assert.NotEmpty(t, result.Key)
		assert.Contains(t, result.URL, result.Key)
		assert.Equal(t, int64(len(content)), result.Size)
	})

	t.Run("creates nested directories", func(t *testing.T) {
		content := []byte("Nested content")
		key := "path/to/nested/file.txt"

		result, err := s.Upload(ctx, ports.UploadRequest{
			Key:         key,
			Content:     bytes.NewReader(content),
			ContentType: "text/plain",
		})
		require.NoError(t, err)

		assert.Equal(t, key, result.Key)

		data, err := os.ReadFile(filepath.Join(tempDir, key))
		require.NoError(t, err)
		assert.Equal(t, content, data)
	})
}

func TestLocalStorageDownload(t *testing.T) {
	tempDir := t.TempDir()

	s, err := storage.NewLocalStorage(storage.LocalConfig{
		BasePath: tempDir,
		BaseURL:  "http://example.com/files",
	})
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("downloads existing file", func(t *testing.T) {
		content := []byte("Downloadable content")
		key := "download-test.txt"

		_, err := s.Upload(ctx, ports.UploadRequest{
			Key:     key,
			Content: bytes.NewReader(content),
		})
		require.NoError(t, err)

		reader, err := s.Download(ctx, key)
		require.NoError(t, err)
		defer reader.Close()

		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, content, data)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		_, err := s.Download(ctx, "non-existent-file.txt")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "file not found")
	})
}

func TestLocalStorageDelete(t *testing.T) {
	tempDir := t.TempDir()

	s, err := storage.NewLocalStorage(storage.LocalConfig{
		BasePath: tempDir,
		BaseURL:  "http://example.com/files",
	})
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("deletes existing file", func(t *testing.T) {
		content := []byte("To be deleted")
		key := "delete-test.txt"

		_, err := s.Upload(ctx, ports.UploadRequest{
			Key:     key,
			Content: bytes.NewReader(content),
		})
		require.NoError(t, err)

		_, err = os.Stat(filepath.Join(tempDir, key))
		require.NoError(t, err)

		err = s.Delete(ctx, key)
		require.NoError(t, err)

		_, err = os.Stat(filepath.Join(tempDir, key))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("does not return error for non-existent file", func(t *testing.T) {
		err := s.Delete(ctx, "already-deleted.txt")
		assert.NoError(t, err)
	})
}

func TestLocalStorageGetURL(t *testing.T) {
	tempDir := t.TempDir()
	baseURL := "http://example.com/files"

	s, err := storage.NewLocalStorage(storage.LocalConfig{
		BasePath: tempDir,
		BaseURL:  baseURL,
	})
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("returns URL for existing file", func(t *testing.T) {
		content := []byte("URL test content")
		key := "url-test.txt"

		_, err := s.Upload(ctx, ports.UploadRequest{
			Key:     key,
			Content: bytes.NewReader(content),
		})
		require.NoError(t, err)

		url, err := s.GetURL(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, baseURL+"/"+key, url)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		_, err := s.GetURL(ctx, "non-existent.txt")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "file not found")
	})
}

func TestLocalStorageClose(t *testing.T) {
	tempDir := t.TempDir()

	s, err := storage.NewLocalStorage(storage.LocalConfig{
		BasePath: tempDir,
		BaseURL:  "http://example.com/files",
	})
	require.NoError(t, err)

	err = s.Close()
	assert.NoError(t, err)
}

func TestDefaultLocalConfig(t *testing.T) {
	cfg := storage.DefaultLocalConfig()

	assert.Equal(t, "./storage", cfg.BasePath)
	assert.Equal(t, "http://localhost:8080/files", cfg.BaseURL)
}
