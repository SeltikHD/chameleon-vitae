package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// contextKey is a type for context keys to avoid collisions.
type contextKey string

// Context keys for storing values in request context.
const (
	// UserContextKey is the key for storing authenticated user info in context.
	UserContextKey contextKey = "user"

	// ClaimsContextKey is the key for storing auth claims in context.
	ClaimsContextKey contextKey = "claims"
)

// AuthenticatedUser represents the authenticated user info stored in context.
type AuthenticatedUser struct {
	ID          string
	FirebaseUID string
	Email       string
}

// GetAuthenticatedUser retrieves the authenticated user from the request context.
func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	user, ok := ctx.Value(UserContextKey).(*AuthenticatedUser)
	return user, ok
}

// GetAuthClaims retrieves the auth claims from the request context.
func GetAuthClaims(ctx context.Context) (*ports.AuthClaims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*ports.AuthClaims)
	return claims, ok
}

// AuthMiddlewareConfig holds configuration for the auth middleware.
type AuthMiddlewareConfig struct {
	AuthProvider ports.AuthProvider
	UserRepo     ports.UserRepository
	SkipPaths    []string
}

// authMiddleware stores the auth dependencies.
type authMiddleware struct {
	authProvider ports.AuthProvider
	userRepo     ports.UserRepository
}

// SetAuthMiddleware sets up the authentication middleware on the router.
func (r *Router) SetAuthMiddleware(authProvider ports.AuthProvider, userRepo ports.UserRepository) {
	r.authMiddleware = &authMiddleware{
		authProvider: authProvider,
		userRepo:     userRepo,
	}
}

// authMiddleware instance stored in router.
type authMiddlewareInstance struct {
	authProvider ports.AuthProvider
	userRepo     ports.UserRepository
}

// AuthMiddleware is the authentication middleware that validates Firebase tokens.
func (r *Router) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Check for auth middleware configuration
		if r.authMiddleware == nil {
			log.Error().Msg("Auth middleware not configured")
			respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Authentication not configured")
			return
		}

		// Extract token from Authorization header
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header")
			return
		}

		// Expect "Bearer <token>" format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization header format")
			return
		}

		token := parts[1]
		if token == "" {
			respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Empty token")
			return
		}

		// Verify the token
		claims, err := r.authMiddleware.authProvider.VerifyToken(req.Context(), token)
		if err != nil {
			log.Debug().Err(err).Msg("Token verification failed")
			respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token")
			return
		}

		// Fetch user from database
		user, err := r.authMiddleware.userRepo.GetByFirebaseUID(req.Context(), claims.UserID)
		if err != nil {
			log.Debug().Err(err).Str("firebase_uid", claims.UserID).Msg("User not found")
			respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not found. Please sync your account first.")
			return
		}

		// Store authenticated user info in context
		authUser := &AuthenticatedUser{
			ID:          user.ID,
			FirebaseUID: user.FirebaseUID,
		}
		if user.Email != nil {
			authUser.Email = *user.Email
		}

		// Add user and claims to context
		ctx := context.WithValue(req.Context(), UserContextKey, authUser)
		ctx = context.WithValue(ctx, ClaimsContextKey, claims)

		// Continue to next handler
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

// ZerologLogger is a middleware that logs requests using zerolog.
func ZerologLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			// Log after request completes
			log.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Dur("latency", time.Since(start)).
				Str("request_id", middleware.GetReqID(r.Context())).
				Msg("HTTP request")
		}()

		next.ServeHTTP(ww, r)
	})
}

// ContentTypeJSON ensures JSON content type for POST/PUT/PATCH requests.
func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only check content type for requests with body
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			contentType := r.Header.Get("Content-Type")
			if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
				respondError(w, http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// CORS returns a middleware that handles CORS headers.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}

			if allowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "300")
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimiter creates a rate limiting middleware.
// Note: This is a placeholder. For production, use a proper rate limiter
// like github.com/go-chi/httprate or implement with Redis.
func RateLimiter(requestsPerMinute int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implement rate limiting with Redis or in-memory store
			// For now, just pass through
			next.ServeHTTP(w, r)
		})
	}
}
