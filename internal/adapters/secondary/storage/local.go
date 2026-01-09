// Package storage provides file storage adapters.
package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// LocalConfig contains configuration for local file storage.
type LocalConfig struct {
	// BasePath is the root directory for file storage.
	BasePath string

	// BaseURL is the base URL for serving files.
	BaseURL string
}

// DefaultLocalConfig returns default local storage configuration.
func DefaultLocalConfig() LocalConfig {
	return LocalConfig{
		BasePath: "./storage",
		BaseURL:  "http://localhost:8080/files",
	}
}

// LocalStorage implements FileStorage using the local filesystem.
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage creates a new local file storage adapter.
func NewLocalStorage(cfg LocalConfig) (*LocalStorage, error) {
	if cfg.BasePath == "" {
		cfg = DefaultLocalConfig()
	}

	// Ensure the base path exists
	if err := os.MkdirAll(cfg.BasePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: cfg.BasePath,
		baseURL:  cfg.BaseURL,
	}, nil
}

// Upload stores a file on the local filesystem.
func (s *LocalStorage) Upload(ctx context.Context, req ports.UploadRequest) (*ports.UploadResult, error) {
	// Generate a unique key if not provided
	key := req.Key
	if key == "" {
		key = uuid.New().String()
	}

	// Create full path
	fullPath := filepath.Join(s.basePath, key)

	// Ensure parent directories exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy content to file
	size, err := io.Copy(file, req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &ports.UploadResult{
		Key:  key,
		URL:  fmt.Sprintf("%s/%s", s.baseURL, key),
		Size: size,
	}, nil
}

// Download retrieves a file from the local filesystem.
func (s *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, key)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", key)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete removes a file from the local filesystem.
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(s.basePath, key)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // File already deleted
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetURL returns the URL for accessing a file.
func (s *LocalStorage) GetURL(ctx context.Context, key string) (string, error) {
	fullPath := filepath.Join(s.basePath, key)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", key)
	}

	return fmt.Sprintf("%s/%s", s.baseURL, key), nil
}

// Close releases any resources held by the storage.
func (s *LocalStorage) Close() error {
	// No resources to release for local storage
	return nil
}

// Ensure LocalStorage implements FileStorage.
var _ ports.FileStorage = (*LocalStorage)(nil)
