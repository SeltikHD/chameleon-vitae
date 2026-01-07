// Package ports defines the interfaces (ports) that adapters must implement.
package ports

import "context"

// AuthClaims represents the claims extracted from an authentication token.
type AuthClaims struct {
	// UserID is the unique identifier from the auth provider (e.g., Firebase UID).
	UserID string

	// Email is the user's email address (may be empty if not provided).
	Email string

	// EmailVerified indicates whether the email has been verified.
	EmailVerified bool

	// Name is the user's display name (may be empty).
	Name string

	// Picture is the URL to the user's profile picture (may be empty).
	Picture string

	// Provider is the authentication provider (e.g., "google.com", "github.com", "password").
	Provider string

	// IssuedAt is the Unix timestamp when the token was issued.
	IssuedAt int64

	// ExpiresAt is the Unix timestamp when the token expires.
	ExpiresAt int64
}

// AuthProvider defines the interface for authentication providers.
// Adapters implementing this interface can verify tokens from external auth services.
type AuthProvider interface {
	// VerifyToken validates an ID token and extracts the claims.
	// Returns an error if the token is invalid, expired, or cannot be verified.
	VerifyToken(ctx context.Context, idToken string) (*AuthClaims, error)

	// Close releases any resources held by the auth provider.
	Close() error
}
