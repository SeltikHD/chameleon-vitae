// Package firebase provides a Firebase Authentication adapter.
package firebase

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

const (
	// googleCertsURL is the URL for Firebase/Google public keys.
	googleCertsURL = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"

	// googleJWKsURL is the alternative JWKs endpoint.
	googleJWKsURL = "https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com"

	// tokenIssuerPrefix is the expected issuer prefix for Firebase tokens.
	tokenIssuerPrefix = "https://securetoken.google.com/"

	// keyRefreshInterval is how often to refresh the public keys.
	keyRefreshInterval = 1 * time.Hour
)

var (
	// ErrInvalidToken is returned when the token is invalid or malformed.
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenExpired is returned when the token has expired.
	ErrTokenExpired = errors.New("token expired")

	// ErrInvalidIssuer is returned when the token issuer doesn't match.
	ErrInvalidIssuer = errors.New("invalid token issuer")

	// ErrInvalidAudience is returned when the token audience doesn't match.
	ErrInvalidAudience = errors.New("invalid token audience")

	// ErrKeyNotFound is returned when the signing key is not found.
	ErrKeyNotFound = errors.New("signing key not found")
)

// JWK represents a JSON Web Key.
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKSet represents a set of JSON Web Keys.
type JWKSet struct {
	Keys []JWK `json:"keys"`
}

// Config holds the configuration for the Firebase auth adapter.
type Config struct {
	// ProjectID is the Firebase project ID (required).
	ProjectID string

	// HTTPClient is an optional HTTP client for fetching public keys.
	// If nil, http.DefaultClient is used.
	HTTPClient *http.Client
}

// Adapter implements the ports.AuthProvider interface for Firebase Authentication.
type Adapter struct {
	projectID  string
	httpClient *http.Client

	mu        sync.RWMutex
	keys      map[string]*rsa.PublicKey
	lastFetch time.Time
}

// New creates a new Firebase authentication adapter.
func New(cfg Config) (*Adapter, error) {
	if cfg.ProjectID == "" {
		return nil, errors.New("firebase: project ID is required")
	}

	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	adapter := &Adapter{
		projectID:  cfg.ProjectID,
		httpClient: client,
		keys:       make(map[string]*rsa.PublicKey),
	}

	// Fetch keys on initialization.
	if err := adapter.refreshKeys(context.Background()); err != nil {
		log.Warn().Err(err).Msg("firebase: failed to fetch initial keys, will retry on first verification")
	}

	return adapter, nil
}

// VerifyToken validates a Firebase ID token and returns the claims.
func (a *Adapter) VerifyToken(ctx context.Context, idToken string) (*ports.AuthClaims, error) {
	// Ensure we have fresh keys.
	if err := a.ensureFreshKeys(ctx); err != nil {
		return nil, fmt.Errorf("failed to fetch signing keys: %w", err)
	}

	// Parse the token without validation first to get the key ID.
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Get the key ID from the header.
	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return nil, fmt.Errorf("%w: missing kid header", ErrInvalidToken)
	}

	// Get the public key for this kid.
	a.mu.RLock()
	pubKey, exists := a.keys[kid]
	a.mu.RUnlock()

	if !exists {
		// Try refreshing keys in case a new key was added.
		if err := a.refreshKeys(ctx); err != nil {
			return nil, fmt.Errorf("failed to refresh signing keys: %w", err)
		}

		a.mu.RLock()
		pubKey, exists = a.keys[kid]
		a.mu.RUnlock()

		if !exists {
			return nil, ErrKeyNotFound
		}
	}

	// Now parse and validate the token.
	expectedIssuer := tokenIssuerPrefix + a.projectID

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(idToken, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "RS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	}, jwt.WithIssuer(expectedIssuer), jwt.WithAudience(a.projectID))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	// Extract claims.
	authClaims := &ports.AuthClaims{}

	if sub, ok := claims["sub"].(string); ok {
		authClaims.UserID = sub
	}
	if email, ok := claims["email"].(string); ok {
		authClaims.Email = email
	}
	if emailVerified, ok := claims["email_verified"].(bool); ok {
		authClaims.EmailVerified = emailVerified
	}
	if name, ok := claims["name"].(string); ok {
		authClaims.Name = name
	}
	if picture, ok := claims["picture"].(string); ok {
		authClaims.Picture = picture
	}

	// Extract provider from firebase claims.
	if firebase, ok := claims["firebase"].(map[string]interface{}); ok {
		if signInProvider, ok := firebase["sign_in_provider"].(string); ok {
			authClaims.Provider = signInProvider
		}
	}

	if iat, ok := claims["iat"].(float64); ok {
		authClaims.IssuedAt = int64(iat)
	}
	if exp, ok := claims["exp"].(float64); ok {
		authClaims.ExpiresAt = int64(exp)
	}

	return authClaims, nil
}

// Close releases resources held by the adapter.
func (a *Adapter) Close() error {
	return nil
}

// ensureFreshKeys checks if keys need to be refreshed and fetches them if necessary.
func (a *Adapter) ensureFreshKeys(ctx context.Context) error {
	a.mu.RLock()
	needsRefresh := time.Since(a.lastFetch) > keyRefreshInterval || len(a.keys) == 0
	a.mu.RUnlock()

	if needsRefresh {
		return a.refreshKeys(ctx)
	}
	return nil
}

// refreshKeys fetches the latest public keys from Google.
func (a *Adapter) refreshKeys(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleJWKsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch keys: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch keys: status %d", resp.StatusCode)
	}

	var jwkSet JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&jwkSet); err != nil {
		return fmt.Errorf("failed to decode JWK set: %w", err)
	}

	newKeys := make(map[string]*rsa.PublicKey)
	for _, jwk := range jwkSet.Keys {
		if jwk.Kty != "RSA" {
			continue
		}

		pubKey, err := jwkToRSAPublicKey(jwk)
		if err != nil {
			log.Warn().Err(err).Str("kid", jwk.Kid).Msg("firebase: failed to parse JWK")
			continue
		}
		newKeys[jwk.Kid] = pubKey
	}

	a.mu.Lock()
	a.keys = newKeys
	a.lastFetch = time.Now()
	a.mu.Unlock()

	log.Debug().Int("key_count", len(newKeys)).Msg("firebase: refreshed public keys")

	return nil
}

// jwkToRSAPublicKey converts a JWK to an RSA public key.
func jwkToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// Decode the modulus.
	nBytes, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(jwk.N, "="))
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}
	n := new(big.Int).SetBytes(nBytes)

	// Decode the exponent.
	eBytes, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(jwk.E, "="))
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	// Convert exponent bytes to int.
	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{N: n, E: e}, nil
}

// Compile-time check that Adapter implements ports.AuthProvider.
var _ ports.AuthProvider = (*Adapter)(nil)
