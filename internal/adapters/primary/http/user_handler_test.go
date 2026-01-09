package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/adapters/primary/http/mocks"
	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

func TestUserHandlerGetMe(t *testing.T) {
	tests := []struct {
		name           string
		setupAuth      func(ctx context.Context) context.Context
		setupMocks     func(userRepo *mocks.InMemoryUserRepository)
		expectedStatus int
		expectedCode   string
		checkResponse  func(t *testing.T, resp UserResponse)
	}{
		{
			name: "success - returns user profile",
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, UserContextKey, &AuthenticatedUser{
					ID:          "user-123",
					FirebaseUID: "firebase-123",
					Email:       "test@example.com",
				})
			},
			setupMocks: func(userRepo *mocks.InMemoryUserRepository) {
				user, _ := domain.NewUser("firebase-123")
				user.ID = "user-123"
				user.SetName("Test User")
				user.SetEmail("test@example.com")
				user.Headline = strPtr("Software Engineer")
				userRepo.Seed(user)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp UserResponse) {
				assert.Equal(t, "user-123", resp.ID)
				assert.Equal(t, "Test User", derefStr(resp.Name))
				assert.Equal(t, "test@example.com", derefStr(resp.Email))
				assert.Equal(t, "Software Engineer", derefStr(resp.Headline))
			},
		},
		{
			name: "error - user not authenticated",
			setupAuth: func(ctx context.Context) context.Context {
				return ctx // No auth user in context
			},
			setupMocks:     func(userRepo *mocks.InMemoryUserRepository) { /* no user seeded */ },
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
		},
		{
			name: "error - user not found",
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, UserContextKey, &AuthenticatedUser{
					ID:          "nonexistent-user",
					FirebaseUID: "firebase-123",
					Email:       "test@example.com",
				})
			},
			setupMocks:     func(userRepo *mocks.InMemoryUserRepository) { /* no user seeded */ },
			expectedStatus: http.StatusNotFound,
			expectedCode:   "USER_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := mocks.NewInMemoryUserRepository()
			authProvider := mocks.NewMockAuthProvider()
			tt.setupMocks(userRepo)

			// Create service and handler
			userService := services.NewUserService(userRepo, authProvider)
			handler := NewUserHandler(userService)

			// Create request
			req := newJSONRequest(t, http.MethodGet, "/v1/me", nil)
			req = req.WithContext(tt.setupAuth(req.Context()))

			// Execute request
			rr := executeRequest(t, req, handler.GetMe)

			// Assert status
			assertStatusCode(t, tt.expectedStatus, rr)

			// Assert response
			if tt.expectedCode != "" {
				var errResp ErrorResponse
				parseJSONResponse(t, rr, &errResp)
				assert.Equal(t, tt.expectedCode, errResp.Error.Code)
			}

			if tt.checkResponse != nil {
				var resp UserResponse
				parseJSONResponse(t, rr, &resp)
				tt.checkResponse(t, resp)
			}
		})
	}
}

func TestUserHandlerUpdateMe(t *testing.T) {
	tests := []struct {
		name           string
		setupAuth      func(ctx context.Context) context.Context
		setupMocks     func(userRepo *mocks.InMemoryUserRepository)
		requestBody    interface{}
		expectedStatus int
		expectedCode   string
		checkResponse  func(t *testing.T, resp UserResponse)
	}{
		{
			name: "success - updates user profile",
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, UserContextKey, &AuthenticatedUser{
					ID:          "user-123",
					FirebaseUID: "firebase-123",
					Email:       "test@example.com",
				})
			},
			setupMocks: func(userRepo *mocks.InMemoryUserRepository) {
				user, _ := domain.NewUser("firebase-123")
				user.ID = "user-123"
				user.SetEmail("test@example.com")
				userRepo.Seed(user)
			},
			requestBody: UpdateUserRequest{
				Name:     strPtr("Updated Name"),
				Headline: strPtr("Senior Engineer"),
				Location: strPtr("San Francisco, CA"),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp UserResponse) {
				assert.Equal(t, "user-123", resp.ID)
				assert.Equal(t, "Updated Name", derefStr(resp.Name))
				assert.Equal(t, "Senior Engineer", derefStr(resp.Headline))
				assert.Equal(t, "San Francisco, CA", derefStr(resp.Location))
			},
		},
		{
			name: "success - partial update",
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, UserContextKey, &AuthenticatedUser{
					ID:          "user-123",
					FirebaseUID: "firebase-123",
					Email:       "test@example.com",
				})
			},
			setupMocks: func(userRepo *mocks.InMemoryUserRepository) {
				user, _ := domain.NewUser("firebase-123")
				user.ID = "user-123"
				user.SetName("Original Name")
				user.Headline = strPtr("Original Headline")
				userRepo.Seed(user)
			},
			requestBody: UpdateUserRequest{
				Headline: strPtr("Updated Headline"),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp UserResponse) {
				assert.Equal(t, "Original Name", derefStr(resp.Name)) // Unchanged
				assert.Equal(t, "Updated Headline", derefStr(resp.Headline))
			},
		},
		{
			name: "error - user not authenticated",
			setupAuth: func(ctx context.Context) context.Context {
				return ctx
			},
			setupMocks: func(userRepo *mocks.InMemoryUserRepository) { /* no user seeded */ },
			requestBody: UpdateUserRequest{
				Name: strPtr("New Name"),
			},
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
		},
		{
			name: "error - user not found",
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, UserContextKey, &AuthenticatedUser{
					ID:          "nonexistent-user",
					FirebaseUID: "firebase-123",
					Email:       "test@example.com",
				})
			},
			setupMocks: func(userRepo *mocks.InMemoryUserRepository) { /* no user seeded */ },
			requestBody: UpdateUserRequest{
				Name: strPtr("New Name"),
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   "USER_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := mocks.NewInMemoryUserRepository()
			authProvider := mocks.NewMockAuthProvider()
			tt.setupMocks(userRepo)

			// Create service and handler
			userService := services.NewUserService(userRepo, authProvider)
			handler := NewUserHandler(userService)

			// Create request
			req := newJSONRequest(t, http.MethodPatch, "/v1/me", tt.requestBody)
			req = req.WithContext(tt.setupAuth(req.Context()))

			// Execute request
			rr := executeRequest(t, req, handler.UpdateMe)

			// Assert status
			assertStatusCode(t, tt.expectedStatus, rr)

			// Assert response
			if tt.expectedCode != "" {
				var errResp ErrorResponse
				parseJSONResponse(t, rr, &errResp)
				assert.Equal(t, tt.expectedCode, errResp.Error.Code)
			}

			if tt.checkResponse != nil {
				var resp UserResponse
				parseJSONResponse(t, rr, &resp)
				tt.checkResponse(t, resp)
			}
		})
	}
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func TestNewUserHandler(t *testing.T) {
	userRepo := mocks.NewInMemoryUserRepository()
	authProvider := mocks.NewMockAuthProvider()
	userService := services.NewUserService(userRepo, authProvider)

	handler := NewUserHandler(userService)
	require.NotNil(t, handler)
	require.NotNil(t, handler.userService)
}
