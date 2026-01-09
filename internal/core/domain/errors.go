// Package domain contains the core business entities and value objects.
package domain

import "errors"

// Domain errors represent business rule violations.
var (
	// User errors.
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidFirebaseUID = errors.New("invalid firebase UID")

	// Experience errors.
	ErrExperienceNotFound    = errors.New("experience not found")
	ErrInvalidExperienceType = errors.New("invalid experience type")
	ErrInvalidDateRange      = errors.New("end date must be after start date")
	ErrInvalidDateFormat     = errors.New("invalid date format, expected YYYY-MM-DD")
	ErrCurrentWithEndDate    = errors.New("current experience cannot have an end date")

	// Bullet errors.
	ErrBulletNotFound     = errors.New("bullet not found")
	ErrEmptyBulletContent = errors.New("bullet content cannot be empty")
	ErrInvalidImpactScore = errors.New("impact score must be between 0 and 100")

	// Skill errors.
	ErrSkillNotFound           = errors.New("skill not found")
	ErrSkillAlreadyExists      = errors.New("skill already exists for this user")
	ErrEmptySkillName          = errors.New("skill name cannot be empty")
	ErrInvalidProficiencyLevel = errors.New("proficiency level must be between 0 and 100")

	// Language errors.
	ErrLanguageNotFound       = errors.New("spoken language not found")
	ErrSpokenLanguageNotFound = errors.New("spoken language not found")
	ErrLanguageAlreadyExists  = errors.New("language already exists for this user")
	ErrInvalidProficiency     = errors.New("invalid language proficiency level")

	// Resume errors.
	ErrResumeNotFound          = errors.New("resume not found")
	ErrInvalidResumeStatus     = errors.New("invalid resume status")
	ErrInvalidMatchScore       = errors.New("match score must be between 0 and 100")
	ErrEmptyJobDescription     = errors.New("job description cannot be empty")
	ErrResumeNotGenerated      = errors.New("resume must be generated before PDF export")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrNoBulletsAvailable      = errors.New("no bullets available for resume generation")
	ErrResumeNotReady          = errors.New("resume is not ready for PDF generation")

	// Validation errors.
	ErrValidation          = errors.New("validation error")
	ErrRequiredField       = errors.New("required field is missing")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidURL          = errors.New("invalid URL format")
	ErrInvalidLanguageCode = errors.New("invalid language code")

	// Authorization errors.
	ErrUnauthorized = errors.New("unauthorized access")
	ErrForbidden    = errors.New("access forbidden")

	// External service errors.
	ErrAIServiceUnavailable  = errors.New("AI service is unavailable")
	ErrPDFServiceUnavailable = errors.New("PDF service is unavailable")
	ErrJobParserUnavailable  = errors.New("job parser service is unavailable")
)

// DomainError wraps a domain error with additional context.
type DomainError struct {
	Err     error
	Message string
	Field   string
}

// Error returns the error message.
func (e *DomainError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}

// Unwrap returns the underlying error.
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error with context.
func NewDomainError(err error, message string) *DomainError {
	return &DomainError{
		Err:     err,
		Message: message,
	}
}

// NewFieldError creates a new domain error for a specific field.
func NewFieldError(err error, field, message string) *DomainError {
	return &DomainError{
		Err:     err,
		Field:   field,
		Message: message,
	}
}

// ValidationErrors collects multiple validation errors.
type ValidationErrors struct {
	Errors []*DomainError
}

// Add adds a validation error.
func (v *ValidationErrors) Add(err *DomainError) {
	v.Errors = append(v.Errors, err)
}

// AddFieldError adds a field validation error.
func (v *ValidationErrors) AddFieldError(field, message string) {
	v.Errors = append(v.Errors, NewFieldError(ErrValidation, field, message))
}

// HasErrors returns true if there are validation errors.
func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// Error returns the combined error message.
func (v *ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return ""
	}
	if len(v.Errors) == 1 {
		return v.Errors[0].Error()
	}
	msg := "multiple validation errors: "
	for i, err := range v.Errors {
		if i > 0 {
			msg += "; "
		}
		msg += err.Error()
	}
	return msg
}

// ToError returns nil if no errors, otherwise returns self.
func (v *ValidationErrors) ToError() error {
	if !v.HasErrors() {
		return nil
	}
	return v
}

// DatabaseError represents a database operation error.
type DatabaseError struct {
	Operation string
	Err       error
}

// Error returns the error message.
func (e *DatabaseError) Error() string {
	return "database error in " + e.Operation + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// NewDatabaseError creates a new database error with operation context.
func NewDatabaseError(operation string, err error) *DatabaseError {
	return &DatabaseError{
		Operation: operation,
		Err:       err,
	}
}
