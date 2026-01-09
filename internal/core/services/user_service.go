// Package services contains the application services (use cases) that orchestrate
// domain logic and coordinate between ports.
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

// UserService handles user-related use cases.
type UserService struct {
	userRepo     ports.UserRepository
	authProvider ports.AuthProvider
}

// NewUserService creates a new UserService with the required dependencies.
func NewUserService(userRepo ports.UserRepository, authProvider ports.AuthProvider) *UserService {
	return &UserService{
		userRepo:     userRepo,
		authProvider: authProvider,
	}
}

// SyncUserRequest contains the parameters for syncing a user from auth provider.
type SyncUserRequest struct {
	IDToken string
}

// SyncUserResponse contains the result of user synchronization.
type SyncUserResponse struct {
	User      *domain.User
	IsNewUser bool
}

// SyncUser synchronizes a user from the authentication provider.
// Creates a new user if they don't exist, or updates their information if they do.
func (s *UserService) SyncUser(ctx context.Context, req SyncUserRequest) (*SyncUserResponse, error) {
	// Verify the ID token and extract claims.
	claims, err := s.authProvider.VerifyToken(ctx, req.IDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	// Try to find existing user.
	existingUser, err := s.userRepo.GetByFirebaseUID(ctx, claims.UserID)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		// Update existing user with latest auth info.
		updated := false

		if claims.Email != "" && (existingUser.Email == nil || *existingUser.Email != claims.Email) {
			existingUser.SetEmail(claims.Email)
			updated = true
		}

		if claims.Name != "" && (existingUser.Name == nil || *existingUser.Name != claims.Name) {
			existingUser.SetName(claims.Name)
			updated = true
		}

		if claims.Picture != "" && (existingUser.PictureURL == nil || *existingUser.PictureURL != claims.Picture) {
			existingUser.PictureURL = &claims.Picture
			existingUser.UpdatedAt = time.Now().UTC()
			updated = true
		}

		if updated {
			if err := s.userRepo.Update(ctx, existingUser); err != nil {
				return nil, fmt.Errorf("failed to update user: %w", err)
			}
		}

		return &SyncUserResponse{
			User:      existingUser,
			IsNewUser: false,
		}, nil
	}

	// Create new user.
	newUser, err := domain.NewUser(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if claims.Email != "" {
		newUser.SetEmail(claims.Email)
	}
	if claims.Name != "" {
		newUser.SetName(claims.Name)
	}
	if claims.Picture != "" {
		newUser.PictureURL = &claims.Picture
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &SyncUserResponse{
		User:      newUser,
		IsNewUser: true,
	}, nil
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetUserByFirebaseUID retrieves a user by their Firebase UID.
func (s *UserService) GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	user, err := s.userRepo.GetByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateProfileRequest contains the parameters for updating a user's profile.
type UpdateProfileRequest struct {
	UserID            string
	Name              *string
	Headline          *string
	Summary           *string
	Location          *string
	Phone             *string
	Website           *string
	LinkedInURL       *string
	GitHubURL         *string
	PortfolioURL      *string
	PreferredLanguage *string
}

// UpdateProfile updates a user's profile information.
func (s *UserService) UpdateProfile(ctx context.Context, req UpdateProfileRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Apply updates.
	if req.Name != nil {
		user.SetName(*req.Name)
	}
	if req.Headline != nil {
		if *req.Headline == "" {
			user.Headline = nil
		} else {
			user.Headline = req.Headline
		}
	}
	if req.Summary != nil {
		if *req.Summary == "" {
			user.Summary = nil
		} else {
			user.Summary = req.Summary
		}
	}
	if req.Location != nil {
		if *req.Location == "" {
			user.Location = nil
		} else {
			user.Location = req.Location
		}
	}
	if req.Phone != nil {
		if *req.Phone == "" {
			user.Phone = nil
		} else {
			user.Phone = req.Phone
		}
	}
	if req.Website != nil {
		if *req.Website == "" {
			user.Website = nil
		} else {
			user.Website = req.Website
		}
	}
	if req.LinkedInURL != nil {
		if *req.LinkedInURL == "" {
			user.LinkedInURL = nil
		} else {
			user.LinkedInURL = req.LinkedInURL
		}
	}
	if req.GitHubURL != nil {
		if *req.GitHubURL == "" {
			user.GitHubURL = nil
		} else {
			user.GitHubURL = req.GitHubURL
		}
	}
	if req.PortfolioURL != nil {
		if *req.PortfolioURL == "" {
			user.PortfolioURL = nil
		} else {
			user.PortfolioURL = req.PortfolioURL
		}
	}
	if req.PreferredLanguage != nil && *req.PreferredLanguage != "" {
		user.PreferredLanguage = *req.PreferredLanguage
	}

	user.UpdatedAt = time.Now().UTC()

	// Validate before saving.
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser removes a user and all their data.
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	// Note: In a real application, you might want to cascade delete
	// experiences, bullets, skills, resumes, etc.
	// This could be handled at the database level with ON DELETE CASCADE
	// or explicitly in this service.
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
