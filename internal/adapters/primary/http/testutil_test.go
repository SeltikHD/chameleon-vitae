package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
)

// testContextKey is a key for test context values.
type testContextKey string

// setupTestContext adds authenticated user to context for testing.
func setupTestContext(userID, firebaseUID, email string) context.Context {
	ctx := context.Background()
	authUser := &AuthenticatedUser{
		ID:          userID,
		FirebaseUID: firebaseUID,
		Email:       email,
	}
	return context.WithValue(ctx, UserContextKey, authUser)
}

// executeRequest executes an HTTP request against a handler and returns the response.
func executeRequest(t *testing.T, req *http.Request, handler http.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// newJSONRequest creates a new HTTP request with JSON body.
func newJSONRequest(t *testing.T, method, path string, body interface{}) *http.Request {
	t.Helper()
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, path, reader)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// newRequestWithChiContext creates a request with Chi URL parameters.
func newRequestWithChiContext(t *testing.T, method, path string, urlParams map[string]string, body interface{}) *http.Request {
	t.Helper()
	req := newJSONRequest(t, method, path, body)

	// Add Chi context with URL parameters
	rctx := chi.NewRouteContext()
	for key, val := range urlParams {
		rctx.URLParams.Add(key, val)
	}
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	return req
}

// parseJSONResponse parses a JSON response body into the given interface.
func parseJSONResponse(t *testing.T, rr *httptest.ResponseRecorder, v interface{}) {
	t.Helper()
	err := json.Unmarshal(rr.Body.Bytes(), v)
	require.NoError(t, err, "Failed to parse JSON response: %s", rr.Body.String())
}

// assertStatusCode asserts the response status code.
func assertStatusCode(t *testing.T, expected int, rr *httptest.ResponseRecorder) {
	t.Helper()
	require.Equal(t, expected, rr.Code, "Response body: %s", rr.Body.String())
}

// assertErrorResponse asserts that the response is an error with the expected code.
func assertErrorResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int, expectedCode string) {
	t.Helper()
	assertStatusCode(t, expectedStatus, rr)

	var errResp ErrorResponse
	parseJSONResponse(t, rr, &errResp)
	require.Equal(t, expectedCode, errResp.Error.Code)
}

// createTestUser creates a test user for testing.
func createTestUser(id string) *domain.User {
	user, _ := domain.NewUser(id)
	user.SetName("Test User")
	user.Headline = strPtr("Software Engineer")
	return user
}

// createTestExperience creates a test experience for testing.
func createTestExperience(id, userID string) *domain.Experience {
	startDate := domain.NewDate(2020, 1, 1)
	exp, _ := domain.NewExperience(
		userID,
		domain.ExperienceTypeWork,
		"Software Engineer",
		"Test Company",
		startDate,
	)
	exp.ID = id
	return exp
}

// createTestBullet creates a test bullet for testing.
func createTestBullet(id, experienceID, _ string) *domain.Bullet {
	bullet, _ := domain.NewBullet(experienceID, "Test bullet content with achievements")
	bullet.ID = id
	return bullet
}

// createTestSkill creates a test skill for testing.
func createTestSkill(id, userID string) *domain.Skill {
	skill, _ := domain.NewSkill(userID, "Go")
	skill.ID = id
	skill.Category = strPtr("Programming Languages")
	return skill
}

// createTestSpokenLanguage creates a test spoken language for testing.
func createTestSpokenLanguage(id, userID string) *domain.SpokenLanguage {
	lang, _ := domain.NewSpokenLanguage(userID, "English", domain.ProficiencyNative)
	lang.ID = id
	return lang
}

// createTestResume creates a test resume for testing.
func createTestResume(id, userID string) *domain.Resume {
	resume, _ := domain.NewResume(userID, "Looking for a talented software engineer...")
	resume.ID = id
	resume.JobTitle = strPtr("Software Engineer at Tech Corp")
	return resume
}
