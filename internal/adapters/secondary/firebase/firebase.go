// Package firebase provides a Firebase Authentication adapter using the official SDK.
package firebase

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"

	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

var (
	// ErrInvalidToken is returned when the token is invalid or malformed.
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenExpired is returned when the token has expired.
	ErrTokenExpired = errors.New("token expired")

	// ErrMissingProjectID is returned when the project ID is not provided.
	ErrMissingProjectID = errors.New("firebase project ID is required")
)

// Config holds the configuration for the Firebase auth adapter.
type Config struct {
	// ProjectID is the Firebase project ID (required).
	ProjectID string

	// CredentialsFile is the path to the service account JSON file.
	// If empty, the SDK will try to use Application Default Credentials.
	CredentialsFile string

	// CredentialsJSON is the service account JSON as a byte array.
	// If provided, takes precedence over CredentialsFile.
	CredentialsJSON []byte
}

// Adapter implements the ports.AuthProvider interface using Firebase Admin SDK.
type Adapter struct {
	app    *firebase.App
	client *auth.Client
}

// New creates a new Firebase authentication adapter using the official SDK.
func New(ctx context.Context, cfg Config) (*Adapter, error) {
	if cfg.ProjectID == "" {
		return nil, ErrMissingProjectID
	}

	var opts []option.ClientOption

	// Configure credentials.
	switch {
	case len(cfg.CredentialsJSON) > 0:
		opts = append(opts, option.WithCredentialsJSON(cfg.CredentialsJSON))
	case cfg.CredentialsFile != "":
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	// If neither is provided, the SDK will attempt to use Application Default Credentials.
	default:
		log.Info().Msg("firebase: using Application Default Credentials")
	}

	// Create Firebase app configuration.
	config := &firebase.Config{
		ProjectID: cfg.ProjectID,
	}

	// Initialize Firebase app.
	app, err := firebase.NewApp(ctx, config, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	// Get auth client.
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get firebase auth client: %w", err)
	}

	log.Info().
		Str("project_id", cfg.ProjectID).
		Msg("firebase: adapter initialized successfully")

	return &Adapter{
		app:    app,
		client: client,
	}, nil
}

// VerifyToken validates a Firebase ID token and returns the claims.
func (a *Adapter) VerifyToken(ctx context.Context, idToken string) (*ports.AuthClaims, error) {
	if idToken == "" {
		return nil, ErrInvalidToken
	}

	// Verify the ID token.
	token, err := a.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		if auth.IsIDTokenExpired(err) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Extract claims.
	claims := &ports.AuthClaims{
		UserID:    token.UID,
		IssuedAt:  token.IssuedAt,
		ExpiresAt: token.Expires,
	}

	// Extract email if present.
	if email, ok := token.Claims["email"].(string); ok {
		claims.Email = email
	}

	// Extract email verified status.
	if emailVerified, ok := token.Claims["email_verified"].(bool); ok {
		claims.EmailVerified = emailVerified
	}

	// Extract name if present.
	if name, ok := token.Claims["name"].(string); ok {
		claims.Name = name
	}

	// Extract picture URL if present.
	if picture, ok := token.Claims["picture"].(string); ok {
		claims.Picture = picture
	}

	// Extract provider from sign_in_provider claim.
	if firebase, ok := token.Claims["firebase"].(map[string]any); ok {
		if provider, ok := firebase["sign_in_provider"].(string); ok {
			claims.Provider = provider
		}
	}

	return claims, nil
}

// GetUser retrieves a user by their Firebase UID.
func (a *Adapter) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	return a.client.GetUser(ctx, uid)
}

// Close releases resources held by the adapter.
func (a *Adapter) Close() error {
	// The Firebase Admin SDK doesn't require explicit cleanup.
	return nil
}

// Compile-time check that Adapter implements ports.AuthProvider.
var _ ports.AuthProvider = (*Adapter)(nil)
