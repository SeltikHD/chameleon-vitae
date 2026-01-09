// Package mocks provides mock implementations of port interfaces for testing.
package mocks

import (
	"context"
	"sync"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// InMemoryUserRepository is an in-memory mock implementation of UserRepository.
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
	// byFirebaseUID maps Firebase UID to user ID for quick lookup.
	byFirebaseUID map[string]string
}

// NewInMemoryUserRepository creates a new in-memory user repository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:         make(map[string]*domain.User),
		byFirebaseUID: make(map[string]string),
	}
}

// Create creates a new user.
func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return domain.ErrUserAlreadyExists
	}

	// Clone user to avoid external mutations
	clone := *user
	r.users[user.ID] = &clone
	r.byFirebaseUID[user.FirebaseUID] = user.ID

	return nil
}

// GetByID retrieves a user by ID.
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	// Return clone to avoid external mutations
	clone := *user
	return &clone, nil
}

// GetByFirebaseUID retrieves a user by Firebase UID.
func (r *InMemoryUserRepository) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.byFirebaseUID[firebaseUID]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	clone := *user
	return &clone, nil
}

// Update updates an existing user.
func (r *InMemoryUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}

	clone := *user
	r.users[user.ID] = &clone
	return nil
}

// Delete removes a user.
func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return domain.ErrUserNotFound
	}

	delete(r.byFirebaseUID, user.FirebaseUID)
	delete(r.users, id)
	return nil
}

// Upsert creates or updates a user.
func (r *InMemoryUserRepository) Upsert(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	clone := *user
	r.users[user.ID] = &clone
	r.byFirebaseUID[user.FirebaseUID] = user.ID
	return nil
}

// Reset clears all data.
func (r *InMemoryUserRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = make(map[string]*domain.User)
	r.byFirebaseUID = make(map[string]string)
}

// Seed adds users for testing.
func (r *InMemoryUserRepository) Seed(users ...*domain.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, user := range users {
		clone := *user
		r.users[user.ID] = &clone
		r.byFirebaseUID[user.FirebaseUID] = user.ID
	}
}

// Verify interface compliance.
var _ ports.UserRepository = (*InMemoryUserRepository)(nil)

// MockAuthProvider is a mock implementation of AuthProvider for testing.
type MockAuthProvider struct {
	mu        sync.RWMutex
	tokens    map[string]*ports.AuthClaims
	verifyErr error
}

// NewMockAuthProvider creates a new mock auth provider.
func NewMockAuthProvider() *MockAuthProvider {
	return &MockAuthProvider{
		tokens: make(map[string]*ports.AuthClaims),
	}
}

// VerifyToken verifies an ID token and returns claims.
func (p *MockAuthProvider) VerifyToken(ctx context.Context, idToken string) (*ports.AuthClaims, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.verifyErr != nil {
		return nil, p.verifyErr
	}

	claims, exists := p.tokens[idToken]
	if !exists {
		return nil, domain.ErrUnauthorized
	}

	return claims, nil
}

// SetVerifyError sets an error to return on VerifyToken.
func (p *MockAuthProvider) SetVerifyError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.verifyErr = err
}

// AddToken adds a token with claims.
func (p *MockAuthProvider) AddToken(token string, claims *ports.AuthClaims) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.tokens[token] = claims
}

// Close closes the mock provider.
func (p *MockAuthProvider) Close() error {
	return nil
}

// Verify interface compliance.
var _ ports.AuthProvider = (*MockAuthProvider)(nil)
