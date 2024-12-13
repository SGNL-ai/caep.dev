package types

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// TransmitterMetadata represents the configuration metadata for a Transmitter
type TransmitterMetadata struct {
	// specVersion identifies the implementer's draft or final specification
	specVersion string `json:"spec_version,omitempty"`

	// issuer is the URL using the HTTPS scheme with no query or fragment component
	issuer       *url.URL `json:"-"`
	issuerString string   `json:"issuer"`

	// jwksUri is the URL of the Transmitter's JSON Web Key Set document
	jwksUri       *url.URL `json:"-"`
	jwksUriString string   `json:"jwks_uri,omitempty"`

	// deliveryMethodsSupported is the list of supported delivery method URIs
	deliveryMethodsSupported []DeliveryMethod `json:"delivery_methods_supported,omitempty"`

	// configurationEndpoint is the URL of the Configuration Endpoint
	configurationEndpoint       *url.URL `json:"-"`
	configurationEndpointString string   `json:"configuration_endpoint,omitempty"`

	// statusEndpoint is the URL of the Status Endpoint
	statusEndpoint       *url.URL `json:"-"`
	statusEndpointString string   `json:"status_endpoint,omitempty"`

	// addSubjectEndpoint is the URL of the Add Subject Endpoint
	addSubjectEndpoint       *url.URL `json:"-"`
	addSubjectEndpointString string   `json:"add_subject_endpoint,omitempty"`

	// removeSubjectEndpoint is the URL of the Remove Subject Endpoint
	removeSubjectEndpoint       *url.URL `json:"-"`
	removeSubjectEndpointString string   `json:"remove_subject_endpoint,omitempty"`

	// verificationEndpoint is the URL of the Verification Endpoint
	verificationEndpoint       *url.URL `json:"-"`
	verificationEndpointString string   `json:"verification_endpoint,omitempty"`

	// criticalSubjectMembers is an array of member names in a Complex Subject which must be interpreted
	criticalSubjectMembers []string `json:"critical_subject_members,omitempty"`

	// authorizationSchemes specifies the supported authorization scheme properties
	authorizationSchemes []AuthorizationScheme `json:"authorization_schemes,omitempty"`

	// defaultSubjects indicates the default behavior of newly created streams
	defaultSubjects string `json:"default_subjects,omitempty"`
}

// AuthorizationScheme represents an authorization scheme supported by the Transmitter
type AuthorizationScheme struct {
	// SpecURN is a URN that describes the specification of the protocol being used
	SpecURN string `json:"spec_urn"`
}

func (m *TransmitterMetadata) UnmarshalJSON(data []byte) error {
	// Create a temporary struct with exported fields for JSON unmarshaling
	var temp struct {
		SpecVersion              string                `json:"spec_version,omitempty"`
		Issuer                   string                `json:"issuer"`
		JWKSUri                  string                `json:"jwks_uri,omitempty"`
		DeliveryMethodsSupported []DeliveryMethod      `json:"delivery_methods_supported,omitempty"`
		ConfigurationEndpoint    string                `json:"configuration_endpoint,omitempty"`
		StatusEndpoint           string                `json:"status_endpoint,omitempty"`
		AddSubjectEndpoint       string                `json:"add_subject_endpoint,omitempty"`
		RemoveSubjectEndpoint    string                `json:"remove_subject_endpoint,omitempty"`
		VerificationEndpoint     string                `json:"verification_endpoint,omitempty"`
		CriticalSubjectMembers   []string              `json:"critical_subject_members,omitempty"`
		AuthorizationSchemes     []AuthorizationScheme `json:"authorization_schemes,omitempty"`
		DefaultSubjects          string                `json:"default_subjects,omitempty"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Copy the values to the actual struct
	m.specVersion = temp.SpecVersion
	m.issuerString = temp.Issuer
	m.jwksUriString = temp.JWKSUri
	m.deliveryMethodsSupported = temp.DeliveryMethodsSupported
	m.configurationEndpointString = temp.ConfigurationEndpoint
	m.statusEndpointString = temp.StatusEndpoint
	m.addSubjectEndpointString = temp.AddSubjectEndpoint
	m.removeSubjectEndpointString = temp.RemoveSubjectEndpoint
	m.verificationEndpointString = temp.VerificationEndpoint
	m.criticalSubjectMembers = temp.CriticalSubjectMembers
	m.authorizationSchemes = temp.AuthorizationSchemes
	m.defaultSubjects = temp.DefaultSubjects

	var err error

	// Parse URLs
	if temp.Issuer != "" {
		if m.issuer, err = url.Parse(temp.Issuer); err != nil {
			return fmt.Errorf("invalid issuer URL: %w", err)
		}
	}

	if temp.JWKSUri != "" {
		if m.jwksUri, err = url.Parse(temp.JWKSUri); err != nil {
			return fmt.Errorf("invalid JWKS URI: %w", err)
		}
	}

	if temp.ConfigurationEndpoint != "" {
		if m.configurationEndpoint, err = url.Parse(temp.ConfigurationEndpoint); err != nil {
			return fmt.Errorf("invalid configuration endpoint URL: %w", err)
		}
	}

	if temp.StatusEndpoint != "" {
		if m.statusEndpoint, err = url.Parse(temp.StatusEndpoint); err != nil {
			return fmt.Errorf("invalid status endpoint URL: %w", err)
		}
	}

	if temp.AddSubjectEndpoint != "" {
		if m.addSubjectEndpoint, err = url.Parse(temp.AddSubjectEndpoint); err != nil {
			return fmt.Errorf("invalid add subject endpoint URL: %w", err)
		}
	}

	if temp.RemoveSubjectEndpoint != "" {
		if m.removeSubjectEndpoint, err = url.Parse(temp.RemoveSubjectEndpoint); err != nil {
			return fmt.Errorf("invalid remove subject endpoint URL: %w", err)
		}
	}

	if temp.VerificationEndpoint != "" {
		if m.verificationEndpoint, err = url.Parse(temp.VerificationEndpoint); err != nil {
			return fmt.Errorf("invalid verification endpoint URL: %w", err)
		}
	}

	return nil
}

func (m *TransmitterMetadata) MarshalJSON() ([]byte, error) {
	// Create a temporary struct with exported fields for JSON marshaling
	temp := struct {
		SpecVersion              string                `json:"spec_version,omitempty"`
		Issuer                   string                `json:"issuer"`
		JWKSUri                  string                `json:"jwks_uri,omitempty"`
		DeliveryMethodsSupported []DeliveryMethod      `json:"delivery_methods_supported,omitempty"`
		ConfigurationEndpoint    string                `json:"configuration_endpoint,omitempty"`
		StatusEndpoint           string                `json:"status_endpoint,omitempty"`
		AddSubjectEndpoint       string                `json:"add_subject_endpoint,omitempty"`
		RemoveSubjectEndpoint    string                `json:"remove_subject_endpoint,omitempty"`
		VerificationEndpoint     string                `json:"verification_endpoint,omitempty"`
		CriticalSubjectMembers   []string              `json:"critical_subject_members,omitempty"`
		AuthorizationSchemes     []AuthorizationScheme `json:"authorization_schemes,omitempty"`
		DefaultSubjects          string                `json:"default_subjects,omitempty"`
	}{
		SpecVersion:              m.specVersion,
		DeliveryMethodsSupported: m.deliveryMethodsSupported,
		CriticalSubjectMembers:   m.criticalSubjectMembers,
		AuthorizationSchemes:     m.authorizationSchemes,
		DefaultSubjects:          m.defaultSubjects,
	}

	// Convert URLs to strings if they exist
	if m.issuer != nil {
		temp.Issuer = m.issuer.String()
	}

	if m.jwksUri != nil {
		temp.JWKSUri = m.jwksUri.String()
	}

	if m.configurationEndpoint != nil {
		temp.ConfigurationEndpoint = m.configurationEndpoint.String()
	}

	if m.statusEndpoint != nil {
		temp.StatusEndpoint = m.statusEndpoint.String()
	}

	if m.addSubjectEndpoint != nil {
		temp.AddSubjectEndpoint = m.addSubjectEndpoint.String()
	}

	if m.removeSubjectEndpoint != nil {
		temp.RemoveSubjectEndpoint = m.removeSubjectEndpoint.String()
	}

	if m.verificationEndpoint != nil {
		temp.VerificationEndpoint = m.verificationEndpoint.String()
	}

	return json.Marshal(temp)
}

func (m *TransmitterMetadata) Validate() error {
	if m.issuer == nil {
		return NewError(
			ErrInvalidTransmitterMetadata,
			"ValidateMetadata",
			"issuer is required",
		)
	}

	if m.configurationEndpoint == nil {
		return NewError(
			ErrInvalidTransmitterMetadata,
			"ValidateMetadata",
			"configuration endpoint is required",
		)
	}

	if len(m.deliveryMethodsSupported) == 0 {
		return NewError(
			ErrInvalidTransmitterMetadata,
			"ValidateMetadata",
			"at least one delivery method must be supported",
		)
	}

	if m.defaultSubjects != "" && !IsValidDefaultSubjects(m.defaultSubjects) {
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
	for _, supported := range m.deliveryMethodsSupported {
		if supported == method {
			return true
		}
	}

	return false
}

func (m *TransmitterMetadata) GetSpecVersion() string {
	return m.specVersion
}

func (m *TransmitterMetadata) GetIssuer() *url.URL {
	if m.issuer == nil {
		return nil
	}

	clone := *m.issuer

	return &clone
}

func (m *TransmitterMetadata) GetJWKSUri() *url.URL {
	if m.jwksUri == nil {
		return nil
	}

	clone := *m.jwksUri

	return &clone
}

func (m *TransmitterMetadata) GetDeliveryMethodsSupported() []DeliveryMethod {
	if m.deliveryMethodsSupported == nil {
		return nil
	}

	methods := make([]DeliveryMethod, len(m.deliveryMethodsSupported))

	copy(methods, m.deliveryMethodsSupported)

	return methods
}

func (m *TransmitterMetadata) GetConfigurationEndpoint() *url.URL {
	if m.configurationEndpoint == nil {
		return nil
	}

	clone := *m.configurationEndpoint

	return &clone
}

func (m *TransmitterMetadata) GetStatusEndpoint() *url.URL {
	if m.statusEndpoint == nil {
		return nil
	}

	clone := *m.statusEndpoint

	return &clone
}

func (m *TransmitterMetadata) GetAddSubjectEndpoint() *url.URL {
	if m.addSubjectEndpoint == nil {
		return nil
	}

	clone := *m.addSubjectEndpoint

	return &clone
}

func (m *TransmitterMetadata) GetRemoveSubjectEndpoint() *url.URL {
	if m.removeSubjectEndpoint == nil {
		return nil
	}

	clone := *m.removeSubjectEndpoint

	return &clone
}

func (m *TransmitterMetadata) GetVerificationEndpoint() *url.URL {
	if m.verificationEndpoint == nil {
		return nil
	}

	clone := *m.verificationEndpoint

	return &clone
}

func (m *TransmitterMetadata) GetCriticalSubjectMembers() []string {
	if m.criticalSubjectMembers == nil {
		return nil
	}

	members := make([]string, len(m.criticalSubjectMembers))

	copy(members, m.criticalSubjectMembers)

	return members
}

func (m *TransmitterMetadata) GetAuthorizationSchemes() []AuthorizationScheme {
	if m.authorizationSchemes == nil {
		return nil
	}

	schemes := make([]AuthorizationScheme, len(m.authorizationSchemes))

	copy(schemes, m.authorizationSchemes)

	return schemes
}

func (m *TransmitterMetadata) GetDefaultSubjects() string {
	return m.defaultSubjects
}
