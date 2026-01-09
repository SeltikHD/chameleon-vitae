package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/adapters/primary/http/mocks"
	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

func TestAuthHandlerSyncUser(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		setupMocks     func(userRepo *mocks.InMemoryUserRepository, authProvider *mocks.MockAuthProvider)
		expectedStatus int
		expectedCode   string
		checkResponse  func(t *testing.T, resp SyncUserResponse)
	}{
		{
			name:       "success - creates new user",
			authHeader: "Bearer valid-token-new-user",
			setupMocks: func(userRepo *mocks.InMemoryUserRepository, authProvider *mocks.MockAuthProvider) {
				authProvider.AddToken("valid-token-new-user", &ports.AuthClaims{
					UserID:        "firebase-new-123",
					Email:         "new@example.com",
					Name:          "New User",
					EmailVerified: true,
				})
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp SyncUserResponse) {
				assert.Equal(t, "firebase-new-123", resp.FirebaseUID)
				assert.Equal(t, "new@example.com", *resp.Email)
				assert.Equal(t, "New User", *resp.Name)
			},
		},
		{
			name:       "success - syncs existing user",
			authHeader: "Bearer valid-token-existing",
			setupMocks: func(userRepo *mocks.InMemoryUserRepository, authProvider *mocks.MockAuthProvider) {
				user, _ := domain.NewUser("firebase-existing-123")
				user.ID = "user-existing-123"
				user.SetEmail("existing@example.com")
				user.SetName("Existing User")
				userRepo.Seed(user)

				authProvider.AddToken("valid-token-existing", &ports.AuthClaims{
					UserID:        "firebase-existing-123",
					Email:         "updated@example.com",
					Name:          "Updated User",
					EmailVerified: true,
				})
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp SyncUserResponse) {
				assert.Equal(t, "firebase-existing-123", resp.FirebaseUID)
				assert.Equal(t, "updated@example.com", *resp.Email)
			},
		},
		{
			name:           "error - missing authorization header",
			authHeader:     "",
			setupMocks:     func(userRepo *mocks.InMemoryUserRepository, authProvider *mocks.MockAuthProvider) { /* no-op */ },
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
		},
		{
			name:           "error - invalid authorization format",
			authHeader:     "InvalidFormat token123",
			setupMocks:     func(userRepo *mocks.InMemoryUserRepository, authProvider *mocks.MockAuthProvider) { /* no-op */ },
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
		},
		{
			name:       "error - invalid token",
			authHeader: "Bearer invalid-token",
			setupMocks: func(userRepo *mocks.InMemoryUserRepository, authProvider *mocks.MockAuthProvider) {
				// Token not added to mock, will fail verification
			},
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := mocks.NewInMemoryUserRepository()
			authProvider := mocks.NewMockAuthProvider()
			tt.setupMocks(userRepo, authProvider)

			// Create service and handler
			userService := services.NewUserService(userRepo, authProvider)
			handler := NewAuthHandler(userService)

			// Create request
			req := newJSONRequest(t, http.MethodPost, "/v1/auth/sync", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			req = req.WithContext(context.Background())

			// Execute request
			rr := executeRequest(t, req, handler.SyncUser)

			// Assert status
			assertStatusCode(t, tt.expectedStatus, rr)

			// Assert error response
			if tt.expectedCode != "" {
				var errResp ErrorResponse
				parseJSONResponse(t, rr, &errResp)
				assert.Equal(t, tt.expectedCode, errResp.Error.Code)
			}

			// Assert success response
			if tt.checkResponse != nil {
				var resp SyncUserResponse
				parseJSONResponse(t, rr, &resp)
				tt.checkResponse(t, resp)
			}
		})
	}
}

func TestNewAuthHandler(t *testing.T) {
	userRepo := mocks.NewInMemoryUserRepository()
	authProvider := mocks.NewMockAuthProvider()
	userService := services.NewUserService(userRepo, authProvider)

	handler := NewAuthHandler(userService)
	require.NotNil(t, handler)
	require.NotNil(t, handler.userService)
}
