package subject

import (
	"encoding/json"
	"fmt"
)

// Format represents the format of a subject identifier
type Format string

const (
	// Standard subject formats
	FormatAccount   Format = "account"
	FormatEmail     Format = "email"
	FormatIssuerSub Format = "iss_sub"
	FormatOpaque    Format = "opaque"
	FormatPhone     Format = "phone_number"
	FormatDID       Format = "did"
	FormatURI       Format = "uri"
	FormatJWTID     Format = "jwt_id"
	FormatSAMLID    Format = "saml_assertion_id"
	FormatComplex   Format = "complex"
)

// ComponentType represents the type of a complex subject component
type ComponentType string

const (
	// Standard component types
	ComponentUser        ComponentType = "user"
	ComponentDevice      ComponentType = "device"
	ComponentSession     ComponentType = "session"
	ComponentApplication ComponentType = "application"
	ComponentTenant      ComponentType = "tenant"
	ComponentOrgUnit     ComponentType = "org_unit"
	ComponentGroup       ComponentType = "group"
)

// Subject is the interface that all subject types must implement
type Subject interface {
	// Format returns the format of the subject identifier
	Format() Format
	// Validate checks if the subject identifier is valid
	Validate() error
	// MarshalJSON implements the json.Marshaler interface
	json.Marshaler
	// UnmarshalJSON implements the json.Unmarshaler interface
	json.Unmarshaler
	// Payload returns the subject's payload as a map[string]interface{}
	Payload() (map[string]interface{}, error)
}

// SimpleSubject represents a subject identified by a single type of identifier
type SimpleSubject interface {
	Subject
}

// ComplexSubject represents a subject identified by multiple component subjects
type ComplexSubject interface {
	Subject

	WithUser(subject Subject) ComplexSubject
	WithDevice(subject Subject) ComplexSubject
	WithSession(subject Subject) ComplexSubject
	WithApplication(subject Subject) ComplexSubject
	WithTenant(subject Subject) ComplexSubject
	WithOrgUnit(subject Subject) ComplexSubject
	WithGroup(subject Subject) ComplexSubject

	// Component getters that return (Subject, bool) to indicate if component exists
	UserComponent() (Subject, bool)
	DeviceComponent() (Subject, bool)
	SessionComponent() (Subject, bool)
	ApplicationComponent() (Subject, bool)
	TenantComponent() (Subject, bool)
	OrgUnitComponent() (Subject, bool)
	GroupComponent() (Subject, bool)
}

// Error types for subject validation
type ErrorCode string

const (
	ErrCodeInvalidFormat    ErrorCode = "invalid_format"
	ErrCodeMissingValue     ErrorCode = "missing_value"
	ErrCodeInvalidValue     ErrorCode = "invalid_value"
	ErrCodeMissingComponent ErrorCode = "missing_component"
)

// SubjectError represents an error that occurred during subject operations
type SubjectError struct {
	Code    ErrorCode
	Message string
	Field   string
}

func (e *SubjectError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Field)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewError creates a new SubjectError
func NewError(code ErrorCode, msg string, field string) error {
	return &SubjectError{
		Code:    code,
		Message: msg,
		Field:   field,
	}
}
