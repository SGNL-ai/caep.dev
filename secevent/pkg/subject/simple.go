package subject

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
)

// baseSimpleSubject provides common functionality for simple subjects
type baseSimpleSubject struct {
	format Format
}

func (s *baseSimpleSubject) Format() Format {
	return s.format
}

// EmailSubject represents a subject identified by an email address
type EmailSubject struct {
	baseSimpleSubject
	email string
}

func NewEmailSubject(email string) (*EmailSubject, error) {
	email = strings.TrimSpace(email)
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, NewError(ErrCodeInvalidValue, "invalid email format", "email")
	}

	return &EmailSubject{
		baseSimpleSubject: baseSimpleSubject{
			format: FormatEmail,
		},
		email: email,
	}, nil
}

func (es *EmailSubject) Email() string {
	return es.email
}

func (es *EmailSubject) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"format": string(es.format),
		"email":  es.email,
	}

	return json.Marshal(m)
}

func (es *EmailSubject) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse subject: %w", err)
	}

	if Format(raw["format"]) != FormatEmail {
		return NewError(ErrCodeInvalidFormat, "invalid format for email subject", "format")
	}

	es.format = FormatEmail
	es.email = strings.TrimSpace(raw["email"])

	return nil
}

func (es *EmailSubject) Validate() error {
	if _, err := mail.ParseAddress(strings.TrimSpace(es.email)); err != nil {
		return NewError(ErrCodeInvalidValue, "invalid email format", "email")
	}

	return nil
}

func (es *EmailSubject) Payload() (map[string]interface{}, error) {
    return map[string]interface{}{
        "format": string(es.format),
        "email":  es.email,
    }, nil
}

// PhoneSubject represents a subject identified by a phone number
type PhoneSubject struct {
	baseSimpleSubject
	phone string
}

func NewPhoneSubject(phone string) (*PhoneSubject, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return nil, NewError(ErrCodeMissingValue, "phone cannot be empty", "phone")
	}

	return &PhoneSubject{
		baseSimpleSubject: baseSimpleSubject{
			format: FormatPhone,
		},
		phone: phone,
	}, nil
}

func (ps *PhoneSubject) Phone() string {
	return ps.phone
}

func (ps *PhoneSubject) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"format": string(ps.format),
		"phone":  ps.phone,
	}

	return json.Marshal(m)
}

func (ps *PhoneSubject) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse subject: %w", err)
	}

	if Format(raw["format"]) != FormatPhone {
		return NewError(ErrCodeInvalidFormat, "invalid format for phone subject", "format")
	}

	ps.format = FormatPhone
	ps.phone = strings.TrimSpace(raw["phone"])

	return nil
}

func (ps *PhoneSubject) Validate() error {
	phone := strings.TrimSpace(ps.phone)
	if phone == "" {
		return NewError(ErrCodeMissingValue, "phone cannot be empty", "phone")
	}
	// TODO: Implement full E.164 validation

	return nil
}

func (ps *PhoneSubject) Payload() (map[string]interface{}, error) {
    return map[string]interface{}{
        "format": string(ps.format),
        "phone":  ps.phone,
    }, nil
}

// IssuerSubSubject represents a subject identified by an issuer and subject pair
type IssuerSubSubject struct {
	baseSimpleSubject
	issuer string
	sub    string
}

func NewIssuerSubSubject(issuer, sub string) (*IssuerSubSubject, error) {
	issuer = strings.TrimSpace(issuer)
	sub = strings.TrimSpace(sub)

	if issuer == "" {
		return nil, NewError(ErrCodeMissingValue, "issuer is required", "issuer")
	}

	if sub == "" {
		return nil, NewError(ErrCodeMissingValue, "subject is required", "sub")
	}

	return &IssuerSubSubject{
		baseSimpleSubject: baseSimpleSubject{
			format: FormatIssuerSub,
		},
		issuer: issuer,
		sub:    sub,
	}, nil
}

func (iss *IssuerSubSubject) Issuer() string {
	return iss.issuer
}

func (iss *IssuerSubSubject) Sub() string {
	return iss.sub
}

func (iss *IssuerSubSubject) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"format": string(iss.format),
		"issuer": iss.issuer,
		"sub":    iss.sub,
	}

	return json.Marshal(m)
}

func (iss *IssuerSubSubject) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse subject: %w", err)
	}

	if Format(raw["format"]) != FormatIssuerSub {
		return NewError(ErrCodeInvalidFormat, "invalid format for issuer_sub subject", "format")
	}

	iss.format = FormatIssuerSub
	iss.issuer = strings.TrimSpace(raw["issuer"])
	iss.sub = strings.TrimSpace(raw["sub"])

	return nil
}

func (iss *IssuerSubSubject) Validate() error {
	if strings.TrimSpace(iss.issuer) == "" {
		return NewError(ErrCodeMissingValue, "issuer is required", "issuer")
	}

	if strings.TrimSpace(iss.sub) == "" {
		return NewError(ErrCodeMissingValue, "subject is required", "sub")
	}

	return nil
}

func (iss *IssuerSubSubject) Payload() (map[string]interface{}, error) {
    return map[string]interface{}{
        "format": string(iss.format),
        "issuer": iss.issuer,
        "sub":    iss.sub,
    }, nil
}

// URISubject represents a subject identified by a URI
type URISubject struct {
	baseSimpleSubject
	uri string
}

func NewURISubject(uri string) (*URISubject, error) {
	uri = strings.TrimSpace(uri)
	if _, err := url.Parse(uri); err != nil {
		return nil, NewError(ErrCodeInvalidValue, "invalid URI format", "uri")
	}

	return &URISubject{
		baseSimpleSubject: baseSimpleSubject{
			format: FormatURI,
		},
		uri: uri,
	}, nil
}

func (us *URISubject) URI() string {
	return us.uri
}

func (us *URISubject) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"format": string(us.format),
		"uri":    us.uri,
	}

	return json.Marshal(m)
}

func (us *URISubject) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse subject: %w", err)
	}

	if Format(raw["format"]) != FormatURI {
		return NewError(ErrCodeInvalidFormat, "invalid format for URI subject", "format")
	}

	us.format = FormatURI
	us.uri = strings.TrimSpace(raw["uri"])

	return nil
}

func (us *URISubject) Validate() error {
	if _, err := url.Parse(strings.TrimSpace(us.uri)); err != nil {
		return NewError(ErrCodeInvalidValue, "invalid URI format", "uri")
	}

	return nil
}

func (us *URISubject) Payload() (map[string]interface{}, error) {
    return map[string]interface{}{
        "format": string(us.format),
        "uri":    us.uri,
    }, nil
}

// OpaqueSubject represents a subject identified by an opaque identifier
type OpaqueSubject struct {
	baseSimpleSubject
	id string
}

func NewOpaqueSubject(id string) (*OpaqueSubject, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, NewError(ErrCodeMissingValue, "identifier is required", "id")
	}

	return &OpaqueSubject{
		baseSimpleSubject: baseSimpleSubject{
			format: FormatOpaque,
		},
		id: id,
	}, nil
}

func (os *OpaqueSubject) ID() string {
	return os.id
}

func (os *OpaqueSubject) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"format": string(os.format),
		"id":     os.id,
	}

	return json.Marshal(m)
}

func (os *OpaqueSubject) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse subject: %w", err)
	}

	if Format(raw["format"]) != FormatOpaque {
		return NewError(ErrCodeInvalidFormat, "invalid format for opaque subject", "format")
	}

	os.format = FormatOpaque
	os.id = strings.TrimSpace(raw["id"])

	return nil
}

func (os *OpaqueSubject) Validate() error {
	if strings.TrimSpace(os.id) == "" {
		return NewError(ErrCodeMissingValue, "identifier is required", "id")
	}

	return nil
}

func (os *OpaqueSubject) Payload() (map[string]interface{}, error) {
    return map[string]interface{}{
        "format": string(os.format),
        "id":     os.id,
    }, nil
}

// AccountSubject represents a subject identified by an acct URI
type AccountSubject struct {
	baseSimpleSubject
	uri string
}

func NewAccountSubject(uri string) (*AccountSubject, error) {
	uri = strings.TrimSpace(uri)
	if !strings.HasPrefix(uri, "acct:") {
		return nil, NewError(ErrCodeInvalidValue, "URI must begin with 'acct:' scheme", "uri")
	}

	if _, err := url.Parse(uri); err != nil {
		return nil, NewError(ErrCodeInvalidValue, "invalid acct URI format", "uri")
	}

	return &AccountSubject{
		baseSimpleSubject: baseSimpleSubject{
			format: FormatAccount,
		},
		uri: uri,
	}, nil
}

func (as *AccountSubject) URI() string {
	return as.uri
}

func (as *AccountSubject) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"format": string(as.format),
		"uri":    as.uri,
	}

	return json.Marshal(m)
}

func (as *AccountSubject) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse subject: %w", err)
	}

	if Format(raw["format"]) != FormatAccount {
		return NewError(ErrCodeInvalidFormat, "invalid format for account subject", "format")
	}

	as.format = FormatAccount
	as.uri = strings.TrimSpace(raw["uri"])

	return nil
}

func (as *AccountSubject) Validate() error {
	uri := strings.TrimSpace(as.uri)
	if !strings.HasPrefix(uri, "acct:") {
		return NewError(ErrCodeInvalidValue, "URI must begin with 'acct:' scheme", "uri")
	}

	if _, err := url.Parse(uri); err != nil {
		return NewError(ErrCodeInvalidValue, "invalid acct URI format", "uri")
	}

	return nil
}

func (as *AccountSubject) Payload() (map[string]interface{}, error) {
    return map[string]interface{}{
        "format": string(as.format),
        "uri":    as.uri,
    }, nil
}

// DIDSubject represents a subject identified by a DID URL
type DIDSubject struct {
	baseSimpleSubject
	url string
}

func NewDIDSubject(url string) (*DIDSubject, error) {
	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "did:") {
		return nil, NewError(ErrCodeInvalidValue, "URL must begin with 'did:' scheme", "url")
	}

	return &DIDSubject{
		baseSimpleSubject: baseSimpleSubject{
			format: FormatDID,
		},
		url: url,
	}, nil
}

func (ds *DIDSubject) URL() string {
	return ds.url
}

func (ds *DIDSubject) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"format": string(ds.format),
		"url":    ds.url,
	}

	return json.Marshal(m)
}

func (ds *DIDSubject) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse subject: %w", err)
	}

	if Format(raw["format"]) != FormatDID {
		return NewError(ErrCodeInvalidFormat, "invalid format for DID subject", "format")
	}

	ds.format = FormatDID
	ds.url = strings.TrimSpace(raw["url"])

	return nil
}

func (ds *DIDSubject) Validate() error {
	url := strings.TrimSpace(ds.url)
	if url == "" {
		return NewError(ErrCodeMissingValue, "URL is required", "url")
	}

	if !strings.HasPrefix(url, "did:") {
		return NewError(ErrCodeInvalidValue, "URL must begin with 'did:' scheme", "url")
	}

	return nil
}

func (ds *DIDSubject) Payload() (map[string]interface{}, error) {
    return map[string]interface{}{
        "format": string(ds.format),
        "url":    ds.url,
    }, nil
}

// JWTIDSubject represents a subject identified by JWT issuer and ID
type JWTIDSubject struct {
	baseSimpleSubject
	iss string
	jti string
}

func NewJWTIDSubject(iss, jti string) (*JWTIDSubject, error) {
	iss = strings.TrimSpace(iss)
	jti = strings.TrimSpace(jti)

	if iss == "" {
		return nil, NewError(ErrCodeMissingValue, "issuer is required", "iss")
	}

	if jti == "" {
		return nil, NewError(ErrCodeMissingValue, "JWT ID is required", "jti")
	}

	return &JWTIDSubject{
		baseSimpleSubject: baseSimpleSubject{
			format: FormatJWTID,
		},
		iss: iss,
		jti: jti,
	}, nil
}

func (js *JWTIDSubject) Issuer() string {
	return js.iss
}

func (js *JWTIDSubject) JWTID() string {
	return js.jti
}

func (js *JWTIDSubject) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"format": string(js.format),
		"iss":    js.iss,
		"jti":    js.jti,
	}

	return json.Marshal(m)
}

func (js *JWTIDSubject) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse subject: %w", err)
	}

	if Format(raw["format"]) != FormatJWTID {
		return NewError(ErrCodeInvalidFormat, "invalid format for JWT ID subject", "format")
	}

	js.format = FormatJWTID
	js.iss = strings.TrimSpace(raw["iss"])
	js.jti = strings.TrimSpace(raw["jti"])

	return nil
}

func (js *JWTIDSubject) Validate() error {
	if strings.TrimSpace(js.iss) == "" {
		return NewError(ErrCodeMissingValue, "issuer is required", "iss")
	}

	if strings.TrimSpace(js.jti) == "" {
		return NewError(ErrCodeMissingValue, "JWT ID is required", "jti")
	}

	return nil
}

func (js *JWTIDSubject) Payload() (map[string]interface{}, error) {
    return map[string]interface{}{
        "format": string(js.format),
        "iss":    js.iss,
        "jti":    js.jti,
    }, nil
}

// SAMLIDSubject represents a subject identified by SAML assertion issuer and ID
type SAMLIDSubject struct {
	baseSimpleSubject
	issuer      string
	assertionID string
}

func NewSAMLIDSubject(issuer, assertionID string) (*SAMLIDSubject, error) {
	issuer = strings.TrimSpace(issuer)
	assertionID = strings.TrimSpace(assertionID)

	if issuer == "" {
		return nil, NewError(ErrCodeMissingValue, "issuer is required", "issuer")
	}

	if assertionID == "" {
		return nil, NewError(ErrCodeMissingValue, "assertion ID is required", "assertion_id")
	}

	return &SAMLIDSubject{
		baseSimpleSubject: baseSimpleSubject{
			format: FormatSAMLID,
		},
		issuer:      issuer,
		assertionID: assertionID,
	}, nil
}

func (ss *SAMLIDSubject) Issuer() string {
	return ss.issuer
}

func (ss *SAMLIDSubject) AssertionID() string {
	return ss.assertionID
}

func (ss *SAMLIDSubject) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"format":       string(ss.format),
		"issuer":       ss.issuer,
		"assertion_id": ss.assertionID,
	}

	return json.Marshal(m)
}

func (ss *SAMLIDSubject) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse subject: %w", err)
	}

	if Format(raw["format"]) != FormatSAMLID {
		return NewError(ErrCodeInvalidFormat, "invalid format for SAML ID subject", "format")
	}

	ss.format = FormatSAMLID
	ss.issuer = strings.TrimSpace(raw["issuer"])
	ss.assertionID = strings.TrimSpace(raw["assertion_id"])

	return nil
}

func (ss *SAMLIDSubject) Validate() error {
	if strings.TrimSpace(ss.issuer) == "" {
		return NewError(ErrCodeMissingValue, "issuer is required", "issuer")
	}

	if strings.TrimSpace(ss.assertionID) == "" {
		return NewError(ErrCodeMissingValue, "assertion ID is required", "assertion_id")
	}

	return nil
}

func (ss *SAMLIDSubject) Payload() (map[string]interface{}, error) {
    return map[string]interface{}{
        "format":       string(ss.format),
        "issuer":       ss.issuer,
        "assertion_id": ss.assertionID,
    }, nil
}
