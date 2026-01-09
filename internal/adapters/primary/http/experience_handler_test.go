package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/adapters/primary/http/mocks"
	"github.com/SeltikHD/chameleon-vitae/internal/core/services"
)

func TestExperienceHandlerList(t *testing.T) {
	tests := []struct {
		name           string
		setupAuth      func(ctx context.Context) context.Context
		setupMocks     func(expRepo *mocks.InMemoryExperienceRepository)
		expectedStatus int
		expectedCode   string
		checkResponse  func(t *testing.T, resp ListExperiencesResponse)
	}{
		{
			name: "success - returns empty list",
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, UserContextKey, &AuthenticatedUser{
					ID:          "user-123",
					FirebaseUID: "firebase-123",
					Email:       "test@example.com",
				})
			},
			setupMocks:     func(expRepo *mocks.InMemoryExperienceRepository) { /* no experiences seeded */ },
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp ListExperiencesResponse) {
				assert.Empty(t, resp.Data)
				assert.Equal(t, 0, resp.Total)
			},
		},
		{
			name: "success - returns experiences",
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, UserContextKey, &AuthenticatedUser{
					ID:          "user-123",
					FirebaseUID: "firebase-123",
					Email:       "test@example.com",
				})
			},
			setupMocks: func(expRepo *mocks.InMemoryExperienceRepository) {
				exp := createTestExperience("exp-1", "user-123")
				expRepo.Seed(exp)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp ListExperiencesResponse) {
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, 1, resp.Total)
				assert.Equal(t, "exp-1", resp.Data[0].ID)
			},
		},
		{
			name: "error - user not authenticated",
			setupAuth: func(ctx context.Context) context.Context {
				return ctx
			},
			setupMocks:     func(expRepo *mocks.InMemoryExperienceRepository) { /* no experiences seeded */ },
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			expRepo := mocks.NewInMemoryExperienceRepository()
			bulletRepo := mocks.NewInMemoryBulletRepository()
			tt.setupMocks(expRepo)

			// Create service and handler
			expService := services.NewExperienceService(expRepo, bulletRepo)
			handler := NewExperienceHandler(expService)

			// Create request
			req := newJSONRequest(t, http.MethodGet, "/v1/experiences", nil)
			req = req.WithContext(tt.setupAuth(req.Context()))

			// Execute request
			rr := executeRequest(t, req, handler.List)

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
				var resp ListExperiencesResponse
				parseJSONResponse(t, rr, &resp)
				tt.checkResponse(t, resp)
			}
		})
	}
}

func TestExperienceHandlerGet(t *testing.T) {
	tests := []struct {
		name           string
		experienceID   string
		setupAuth      func(ctx context.Context) context.Context
		setupMocks     func(expRepo *mocks.InMemoryExperienceRepository)
		expectedStatus int
		expectedCode   string
		checkResponse  func(t *testing.T, resp ExperienceResponse)
	}{
		{
			name:         "success - returns experience",
			experienceID: "exp-1",
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, UserContextKey, &AuthenticatedUser{
					ID:          "user-123",
					FirebaseUID: "firebase-123",
					Email:       "test@example.com",
				})
			},
			setupMocks: func(expRepo *mocks.InMemoryExperienceRepository) {
				exp := createTestExperience("exp-1", "user-123")
				expRepo.Seed(exp)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp ExperienceResponse) {
				assert.Equal(t, "exp-1", resp.ID)
				assert.Equal(t, "work", resp.Type)
			},
		},
		{
			name:         "error - experience not found",
			experienceID: "non-existent",
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, UserContextKey, &AuthenticatedUser{
					ID:          "user-123",
					FirebaseUID: "firebase-123",
					Email:       "test@example.com",
				})
			},
			setupMocks:     func(expRepo *mocks.InMemoryExperienceRepository) { /* no experiences seeded */ },
			expectedStatus: http.StatusNotFound,
			expectedCode:   "EXPERIENCE_NOT_FOUND",
		},
		{
			name:         "error - unauthorized access (returns not found for security)",
			experienceID: "exp-1",
			setupAuth: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, UserContextKey, &AuthenticatedUser{
					ID:          "user-different",
					FirebaseUID: "firebase-different",
					Email:       "different@example.com",
				})
			},
			setupMocks: func(expRepo *mocks.InMemoryExperienceRepository) {
				exp := createTestExperience("exp-1", "user-123") // Different user owns this
				expRepo.Seed(exp)
			},
			// Returns 404 instead of 403 to not reveal existence of other users' experiences
			expectedStatus: http.StatusNotFound,
			expectedCode:   "EXPERIENCE_NOT_FOUND",
		},
		{
			name:         "error - user not authenticated",
			experienceID: "exp-1",
			setupAuth: func(ctx context.Context) context.Context {
				return ctx
			},
			setupMocks:     func(expRepo *mocks.InMemoryExperienceRepository) { /* no experiences seeded */ },
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			expRepo := mocks.NewInMemoryExperienceRepository()
			bulletRepo := mocks.NewInMemoryBulletRepository()
			tt.setupMocks(expRepo)

			// Create service and handler
			expService := services.NewExperienceService(expRepo, bulletRepo)
			handler := NewExperienceHandler(expService)

			// Create base context with auth
			authCtx := tt.setupAuth(context.Background())

			// Create request with URL params and auth context
			req := newRequestWithChiContext(t, http.MethodGet, "/v1/experiences/"+tt.experienceID, map[string]string{
				"experienceID": tt.experienceID,
			}, nil)

			// Merge auth context into Chi context
			mergedCtx := context.WithValue(req.Context(), UserContextKey, authCtx.Value(UserContextKey))
			req = req.WithContext(mergedCtx)

			// Execute request
			rr := executeRequest(t, req, handler.Get)

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
				var resp ExperienceResponse
				parseJSONResponse(t, rr, &resp)
				tt.checkResponse(t, resp)
			}
		})
	}
}

func TestNewExperienceHandler(t *testing.T) {
	expRepo := mocks.NewInMemoryExperienceRepository()
	bulletRepo := mocks.NewInMemoryBulletRepository()
	expService := services.NewExperienceService(expRepo, bulletRepo)

	handler := NewExperienceHandler(expService)
	require.NotNil(t, handler)
	require.NotNil(t, handler.experienceService)
}
