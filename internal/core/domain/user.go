// Package domain contains the core business entities and value objects.
package domain

import "time"

// User represents a system user linked to Firebase authentication.
type User struct {
	ID                string    `json:"id"`
	FirebaseUID       string    `json:"firebase_uid"`
	PictureURL        *string   `json:"picture_url,omitempty"`
	Email             *string   `json:"email,omitempty"`
	Name              *string   `json:"name,omitempty"`
	Headline          *string   `json:"headline,omitempty"`
	Summary           *string   `json:"summary,omitempty"`
	Location          *string   `json:"location,omitempty"`
	Phone             *string   `json:"phone,omitempty"`
	Website           *string   `json:"website,omitempty"`
	LinkedInURL       *string   `json:"linkedin_url,omitempty"`
	GitHubURL         *string   `json:"github_url,omitempty"`
	PortfolioURL      *string   `json:"portfolio_url,omitempty"`
	PreferredLanguage string    `json:"preferred_language"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// NewUser creates a new user with required fields.
func NewUser(firebaseUID string) (*User, error) {
	if firebaseUID == "" {
		return nil, ErrInvalidFirebaseUID
	}

	now := time.Now().UTC()
	return &User{
		FirebaseUID:       firebaseUID,
		PreferredLanguage: "en",
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

// Validate validates the user entity.
func (u *User) Validate() error {
	v := &ValidationErrors{}

	if u.FirebaseUID == "" {
		v.AddFieldError("firebase_uid", "firebase UID is required")
	}

	if u.PreferredLanguage != "" && u.PreferredLanguage != "en" && u.PreferredLanguage != "pt-br" {
		v.AddFieldError("preferred_language", "must be 'en' or 'pt-br'")
	}

	return v.ToError()
}

// SetEmail sets the user's email.
func (u *User) SetEmail(email string) {
	if email == "" {
		u.Email = nil
	} else {
		u.Email = &email
	}
	u.UpdatedAt = time.Now().UTC()
}

// SetName sets the user's name.
func (u *User) SetName(name string) {
	if name == "" {
		u.Name = nil
	} else {
		u.Name = &name
	}
	u.UpdatedAt = time.Now().UTC()
}

// GetDisplayName returns the user's display name (name or email fallback).
func (u *User) GetDisplayName() string {
	if u.Name != nil && *u.Name != "" {
		return *u.Name
	}
	if u.Email != nil && *u.Email != "" {
		return *u.Email
	}
	return "Anonymous User"
}
