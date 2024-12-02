package subject

import (
	"encoding/json"
	"fmt"
)

// ComplexSubjectImpl implements the ComplexSubject interface
type ComplexSubjectImpl struct {
	format      Format
	user        Subject
	device      Subject
	session     Subject
	application Subject
	tenant      Subject
	orgUnit     Subject
	group       Subject
}

// NewComplexSubject creates a new complex subject
func NewComplexSubject() *ComplexSubjectImpl {
	return &ComplexSubjectImpl{
		format: FormatComplex,
	}
}

// Format returns the format of the subject identifier
func (s *ComplexSubjectImpl) Format() Format {
	return s.format
}

// WithUser adds a User component
func (s *ComplexSubjectImpl) WithUser(subject Subject) ComplexSubject {
	if subject != nil {
		s.user = subject
	}

	return s
}

// WithDevice adds a Device component
func (s *ComplexSubjectImpl) WithDevice(subject Subject) ComplexSubject {
	if subject != nil {
		s.device = subject
	}

	return s
}

// WithSession adds a Session component
func (s *ComplexSubjectImpl) WithSession(subject Subject) ComplexSubject {
	if subject != nil {
		s.session = subject
	}

	return s
}

// WithApplication adds an Application component
func (s *ComplexSubjectImpl) WithApplication(subject Subject) ComplexSubject {
	if subject != nil {
		s.application = subject
	}

	return s
}

// WithTenant adds a Tenant component
func (s *ComplexSubjectImpl) WithTenant(subject Subject) ComplexSubject {
	if subject != nil {
		s.tenant = subject
	}

	return s
}

// WithOrgUnit adds an OrgUnit component
func (s *ComplexSubjectImpl) WithOrgUnit(subject Subject) ComplexSubject {
	if subject != nil {
		s.orgUnit = subject
	}

	return s
}

// WithGroup adds a Group component
func (s *ComplexSubjectImpl) WithGroup(subject Subject) ComplexSubject {
	if subject != nil {
		s.group = subject
	}

	return s
}

// Component getters that return (Subject, bool) to indicate if component exists
func (s *ComplexSubjectImpl) UserComponent() (Subject, bool) {
	return s.user, s.user != nil
}

func (s *ComplexSubjectImpl) DeviceComponent() (Subject, bool) {
	return s.device, s.device != nil
}

func (s *ComplexSubjectImpl) SessionComponent() (Subject, bool) {
	return s.session, s.session != nil
}

func (s *ComplexSubjectImpl) ApplicationComponent() (Subject, bool) {
	return s.application, s.application != nil
}

func (s *ComplexSubjectImpl) TenantComponent() (Subject, bool) {
	return s.tenant, s.tenant != nil
}

func (s *ComplexSubjectImpl) OrgUnitComponent() (Subject, bool) {
	return s.orgUnit, s.orgUnit != nil
}

func (s *ComplexSubjectImpl) GroupComponent() (Subject, bool) {
	return s.group, s.group != nil
}

// Validate ensures the complex subject and all its components are valid
func (s *ComplexSubjectImpl) Validate() error {
	if s.format != FormatComplex {
		return NewError(ErrCodeInvalidFormat, "invalid format for complex subject", "format")
	}

	if s.user == nil && s.device == nil && s.session == nil &&
		s.application == nil && s.tenant == nil && s.orgUnit == nil && s.group == nil {

		return NewError(ErrCodeMissingComponent, "complex subject must have at least one component", "")
	}

	// Validate each component if it exists
	if s.user != nil {
		if err := s.user.Validate(); err != nil {
			return NewError(ErrCodeInvalidValue, fmt.Sprintf("invalid user component: %v", err), "user")
		}
	}

	if s.device != nil {
		if err := s.device.Validate(); err != nil {
			return NewError(ErrCodeInvalidValue, fmt.Sprintf("invalid device component: %v", err), "device")
		}
	}

	if s.session != nil {
		if err := s.session.Validate(); err != nil {
			return NewError(ErrCodeInvalidValue, fmt.Sprintf("invalid session component: %v", err), "session")
		}
	}

	if s.application != nil {
		if err := s.application.Validate(); err != nil {
			return NewError(ErrCodeInvalidValue, fmt.Sprintf("invalid application component: %v", err), "application")
		}
	}

	if s.tenant != nil {
		if err := s.tenant.Validate(); err != nil {
			return NewError(ErrCodeInvalidValue, fmt.Sprintf("invalid tenant component: %v", err), "tenant")
		}
	}

	if s.orgUnit != nil {
		if err := s.orgUnit.Validate(); err != nil {
			return NewError(ErrCodeInvalidValue, fmt.Sprintf("invalid org unit component: %v", err), "org_unit")
		}
	}

	if s.group != nil {
		if err := s.group.Validate(); err != nil {
			return NewError(ErrCodeInvalidValue, fmt.Sprintf("invalid group component: %v", err), "group")
		}
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (s *ComplexSubjectImpl) MarshalJSON() ([]byte, error) {
	// Manually construct the JSON object since fields are unexported
	data := map[string]interface{}{
		"format": s.format,
	}

	// Include components if they are not nil
	if s.user != nil {
		data["user"] = s.user
	}

	if s.device != nil {
		data["device"] = s.device
	}

	if s.session != nil {
		data["session"] = s.session
	}

	if s.application != nil {
		data["application"] = s.application
	}

	if s.tenant != nil {
		data["tenant"] = s.tenant
	}

	if s.orgUnit != nil {
		data["org_unit"] = s.orgUnit
	}

	if s.group != nil {
		data["group"] = s.group
	}

	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (s *ComplexSubjectImpl) UnmarshalJSON(data []byte) error {
	// Create a temporary structure to hold the data
	temp := struct {
		Format      Format          `json:"format"`
		User        json.RawMessage `json:"user,omitempty"`
		Device      json.RawMessage `json:"device,omitempty"`
		Session     json.RawMessage `json:"session,omitempty"`
		Application json.RawMessage `json:"application,omitempty"`
		Tenant      json.RawMessage `json:"tenant,omitempty"`
		OrgUnit     json.RawMessage `json:"org_unit,omitempty"`
		Group       json.RawMessage `json:"group,omitempty"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.Format != FormatComplex {
		return NewError(ErrCodeInvalidFormat, "invalid format for complex subject", "format")
	}

	s.format = temp.Format

	// Parse each component if present
	if len(temp.User) > 0 {
		subject, err := ParseSubject(temp.User)
		if err != nil {
			return fmt.Errorf("failed to parse user component: %w", err)
		}

		s.user = subject
	}

	if len(temp.Device) > 0 {
		subject, err := ParseSubject(temp.Device)
		if err != nil {
			return fmt.Errorf("failed to parse device component: %w", err)
		}

		s.device = subject
	}

	if len(temp.Session) > 0 {
		subject, err := ParseSubject(temp.Session)
		if err != nil {
			return fmt.Errorf("failed to parse session component: %w", err)
		}

		s.session = subject
	}

	if len(temp.Application) > 0 {
		subject, err := ParseSubject(temp.Application)
		if err != nil {
			return fmt.Errorf("failed to parse application component: %w", err)
		}

		s.application = subject
	}

	if len(temp.Tenant) > 0 {
		subject, err := ParseSubject(temp.Tenant)
		if err != nil {
			return fmt.Errorf("failed to parse tenant component: %w", err)
		}

		s.tenant = subject
	}

	if len(temp.OrgUnit) > 0 {
		subject, err := ParseSubject(temp.OrgUnit)
		if err != nil {
			return fmt.Errorf("failed to parse org unit component: %w", err)
		}

		s.orgUnit = subject
	}

	if len(temp.Group) > 0 {
		subject, err := ParseSubject(temp.Group)
		if err != nil {
			return fmt.Errorf("failed to parse group component: %w", err)
		}

		s.group = subject
	}

	return nil
}

func (s *ComplexSubjectImpl) Payload() (map[string]interface{}, error) {
    payload := map[string]interface{}{
        "format": string(s.format),
    }

    // Include components if they are not nil
    if s.user != nil {
        userPayload, err := s.user.Payload()
        if err != nil {
            return nil, fmt.Errorf("failed to get user payload: %w", err)
        }
        payload["user"] = userPayload
    }

    if s.device != nil {
        devicePayload, err := s.device.Payload()
        if err != nil {
            return nil, fmt.Errorf("failed to get device payload: %w", err)
        }
        payload["device"] = devicePayload
    }

    if s.session != nil {
        sessionPayload, err := s.session.Payload()
        if err != nil {
            return nil, fmt.Errorf("failed to get session payload: %w", err)
        }
        payload["session"] = sessionPayload
    }

    if s.application != nil {
        appPayload, err := s.application.Payload()
        if err != nil {
            return nil, fmt.Errorf("failed to get application payload: %w", err)
        }
        payload["application"] = appPayload
    }

    if s.tenant != nil {
        tenantPayload, err := s.tenant.Payload()
        if err != nil {
            return nil, fmt.Errorf("failed to get tenant payload: %w", err)
        }
        payload["tenant"] = tenantPayload
    }

    if s.orgUnit != nil {
        orgUnitPayload, err := s.orgUnit.Payload()
        if err != nil {
            return nil, fmt.Errorf("failed to get org unit payload: %w", err)
        }
        payload["org_unit"] = orgUnitPayload
    }

    if s.group != nil {
        groupPayload, err := s.group.Payload()
        if err != nil {
            return nil, fmt.Errorf("failed to get group payload: %w", err)
        }
        payload["group"] = groupPayload
    }

    return payload, nil
}