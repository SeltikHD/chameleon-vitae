package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/adapters/primary/http/mocks"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

func TestNewRouter(t *testing.T) {
	// Create all mock repositories
	userRepo := mocks.NewInMemoryUserRepository()
	authProvider := mocks.NewMockAuthProvider()
	expRepo := mocks.NewInMemoryExperienceRepository()
	bulletRepo := mocks.NewInMemoryBulletRepository()

	// Create services
	userService := services.NewUserService(userRepo, authProvider)
	expService := services.NewExperienceService(expRepo, bulletRepo)
	// Note: Some services are nil for this simple test

	config := DefaultRouterConfig()
	svc := Services{
		UserService:       userService,
		ExperienceService: expService,
		// Other services would be needed for full testing
	}

	router := NewRouter(config, svc)
	require.NotNil(t, router)
}

func TestRouterHealthEndpoint(t *testing.T) {
	// Create all mock repositories
	userRepo := mocks.NewInMemoryUserRepository()
	authProvider := mocks.NewMockAuthProvider()
	expRepo := mocks.NewInMemoryExperienceRepository()
	bulletRepo := mocks.NewInMemoryBulletRepository()

	// Create services
	userService := services.NewUserService(userRepo, authProvider)
	expService := services.NewExperienceService(expRepo, bulletRepo)

	config := DefaultRouterConfig()
	config.EnableSwagger = false // Disable for testing

	svc := Services{
		UserService:       userService,
		ExperienceService: expService,
	}

	router := NewRouter(config, svc)

	// Test health endpoint
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp HealthResponse
	parseJSONResponse(t, rr, &resp)
	assert.Equal(t, "healthy", resp.Status)
	assert.Equal(t, "chameleon-vitae", resp.Service)
}

func TestRouterNotFoundHandler(t *testing.T) {
	// Create all mock repositories
	userRepo := mocks.NewInMemoryUserRepository()
	authProvider := mocks.NewMockAuthProvider()
	expRepo := mocks.NewInMemoryExperienceRepository()
	bulletRepo := mocks.NewInMemoryBulletRepository()

	// Create services
	userService := services.NewUserService(userRepo, authProvider)
	expService := services.NewExperienceService(expRepo, bulletRepo)

	config := DefaultRouterConfig()
	config.EnableSwagger = false

	svc := Services{
		UserService:       userService,
		ExperienceService: expService,
	}

	router := NewRouter(config, svc)

	// Test non-existent endpoint
	req, err := http.NewRequest(http.MethodGet, "/non-existent-path", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var errResp ErrorResponse
	parseJSONResponse(t, rr, &errResp)
	assert.Equal(t, "RESOURCE_NOT_FOUND", errResp.Error.Code)
}

func TestDefaultRouterConfig(t *testing.T) {
	cfg := DefaultRouterConfig()

	assert.True(t, cfg.EnableSwagger)
	assert.False(t, cfg.EnableProfiling)
	assert.Equal(t, int64(10*1024*1024), cfg.MaxRequestSize)
	assert.NotEmpty(t, cfg.AllowedOrigins)
}
