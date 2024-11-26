package types

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// TransmitterMetadata represents the configuration metadata for a Transmitter
type TransmitterMetadata struct {
	// SpecVersion identifies the implementer's draft or final specification
	SpecVersion string `json:"spec_version,omitempty"`

	// Issuer is the URL using the HTTPS scheme with no query or fragment component
	Issuer       *url.URL `json:"-"`
	IssuerString string   `json:"issuer"`

	// JWKSUri is the URL of the Transmitter's JSON Web Key Set document
	JWKSUri       *url.URL `json:"-"`
	JWKSUriString string   `json:"jwks_uri,omitempty"`

	// DeliveryMethodsSupported is the list of supported delivery method URIs
	DeliveryMethodsSupported []DeliveryMethod `json:"delivery_methods_supported,omitempty"`

	// ConfigurationEndpoint is the URL of the Configuration Endpoint
	ConfigurationEndpoint       *url.URL `json:"-"`
	ConfigurationEndpointString string   `json:"configuration_endpoint,omitempty"`

	// StatusEndpoint is the URL of the Status Endpoint
	StatusEndpoint       *url.URL `json:"-"`
	StatusEndpointString string   `json:"status_endpoint,omitempty"`

	// AddSubjectEndpoint is the URL of the Add Subject Endpoint
	AddSubjectEndpoint       *url.URL `json:"-"`
	AddSubjectEndpointString string   `json:"add_subject_endpoint,omitempty"`

	// RemoveSubjectEndpoint is the URL of the Remove Subject Endpoint
	RemoveSubjectEndpoint       *url.URL `json:"-"`
	RemoveSubjectEndpointString string   `json:"remove_subject_endpoint,omitempty"`

	// VerificationEndpoint is the URL of the Verification Endpoint
	VerificationEndpoint       *url.URL `json:"-"`
	VerificationEndpointString string   `json:"verification_endpoint,omitempty"`

	// CriticalSubjectMembers is an array of member names in a Complex Subject which must be interpreted
	CriticalSubjectMembers []string `json:"critical_subject_members,omitempty"`

	// AuthorizationSchemes specifies the supported authorization scheme properties
	AuthorizationSchemes []AuthorizationScheme `json:"authorization_schemes,omitempty"`

	// DefaultSubjects indicates the default behavior of newly created streams
	DefaultSubjects string `json:"default_subjects,omitempty"`
}

// AuthorizationScheme represents an authorization scheme supported by the Transmitter
type AuthorizationScheme struct {
	// SpecURN is a URN that describes the specification of the protocol being used
	SpecURN string `json:"spec_urn"`
}

func (m *TransmitterMetadata) UnmarshalJSON(data []byte) error {
	type TempMetadata TransmitterMetadata

	var temp TempMetadata
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*m = TransmitterMetadata(temp)

	var err error

	if temp.IssuerString != "" {
		if m.Issuer, err = url.Parse(temp.IssuerString); err != nil {
			return fmt.Errorf("invalid issuer URL: %w", err)
		}
	}

	if temp.JWKSUriString != "" {
		if m.JWKSUri, err = url.Parse(temp.JWKSUriString); err != nil {
			return fmt.Errorf("invalid JWKS URI: %w", err)
		}
	}

	if temp.ConfigurationEndpointString != "" {
		if m.ConfigurationEndpoint, err = url.Parse(temp.ConfigurationEndpointString); err != nil {
			return fmt.Errorf("invalid configuration endpoint URL: %w", err)
		}
	}

	if temp.StatusEndpointString != "" {
		if m.StatusEndpoint, err = url.Parse(temp.StatusEndpointString); err != nil {
			return fmt.Errorf("invalid status endpoint URL: %w", err)
		}
	}

	if temp.AddSubjectEndpointString != "" {
		if m.AddSubjectEndpoint, err = url.Parse(temp.AddSubjectEndpointString); err != nil {
			return fmt.Errorf("invalid add subject endpoint URL: %w", err)
		}
	}

	if temp.RemoveSubjectEndpointString != "" {
		if m.RemoveSubjectEndpoint, err = url.Parse(temp.RemoveSubjectEndpointString); err != nil {
			return fmt.Errorf("invalid remove subject endpoint URL: %w", err)
		}
	}

	if temp.VerificationEndpointString != "" {
		if m.VerificationEndpoint, err = url.Parse(temp.VerificationEndpointString); err != nil {
			return fmt.Errorf("invalid verification endpoint URL: %w", err)
		}
	}

	return nil
}

func (m *TransmitterMetadata) MarshalJSON() ([]byte, error) {
	type TempMetadata TransmitterMetadata

	temp := TempMetadata(*m)

	if m.Issuer != nil {
		temp.IssuerString = m.Issuer.String()
	}

	if m.JWKSUri != nil {
		temp.JWKSUriString = m.JWKSUri.String()
	}

	if m.ConfigurationEndpoint != nil {
		temp.ConfigurationEndpointString = m.ConfigurationEndpoint.String()
	}

	if m.StatusEndpoint != nil {
		temp.StatusEndpointString = m.StatusEndpoint.String()
	}

	if m.AddSubjectEndpoint != nil {
		temp.AddSubjectEndpointString = m.AddSubjectEndpoint.String()
	}

	if m.RemoveSubjectEndpoint != nil {
		temp.RemoveSubjectEndpointString = m.RemoveSubjectEndpoint.String()
	}

	if m.VerificationEndpoint != nil {
		temp.VerificationEndpointString = m.VerificationEndpoint.String()
	}

	return json.Marshal(temp)
}

func (m *TransmitterMetadata) Validate() error {
	if m.Issuer == nil {
		return NewError(
			ErrInvalidTransmitterMetadata,
			"ValidateMetadata",
			"issuer is required",
		)
	}

	if m.ConfigurationEndpoint == nil {
		return NewError(
			ErrInvalidTransmitterMetadata,
			"ValidateMetadata",
			"configuration endpoint is required",
		)
	}

	if len(m.DeliveryMethodsSupported) == 0 {
		return NewError(
			ErrInvalidTransmitterMetadata,
			"ValidateMetadata",
			"at least one delivery method must be supported",
		)
	}

	if m.DefaultSubjects != "" && !IsValidDefaultSubjects(m.DefaultSubjects) {
		return NewError(
			ErrInvalidTransmitterMetadata,
			"ValidateMetadata",
			"default_subjects must be either 'ALL' or 'NONE'",
		)
	}

	return nil
}

func IsValidDefaultSubjects(value string) bool {
	return value == DefaultSubjectsAll || value == DefaultSubjectsNone
}

func (m *TransmitterMetadata) SupportsDeliveryMethod(method DeliveryMethod) bool {
	for _, supported := range m.DeliveryMethodsSupported {
		if supported == method {
			return true
		}
	}

	return false
}
