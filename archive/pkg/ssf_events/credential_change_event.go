package ssf_events

type CredentialType string

const (
	Password              CredentialType = "password"
	Pin                   CredentialType = "pin"
	X509                  CredentialType = "x509"
	Fido2_platform        CredentialType = "fido2_platform"
	Fido2_roaming         CredentialType = "fido2_roaming"
	Fido_u2f              CredentialType = "fido_u2f"
	Verifiable_credential CredentialType = "verifiable_credential"
	Phone_voice           CredentialType = "phone_voice"
	Phone_sms             CredentialType = "phone_sms"
	App                   CredentialType = "app"
)

type ChangeType string

const (
	Create  ChangeType = "create"
	Revoked ChangeType = "revoked"
	Update  ChangeType = "update"
	Delete  ChangeType = "delete"
)

var CredentialTypesMap = map[string]CredentialType{
	"password":              Password,
	"pin":                   Pin,
	"x509":                  X509,
	"fido2-platform":        Fido2_platform,
	"fido2-roaming":         Fido2_roaming,
	"fido-u2f":              Fido_u2f,
	"verifiable-credential": Verifiable_credential,
	"phone-voice":           Phone_voice,
	"phone-sms":             Phone_sms,
	"app":                   App,
}

var ChangeTypesMap = map[string]ChangeType{
	"create": Create,
	"revoke": Revoked,
	"update": Update,
	"delete": Delete,
}

// The credential change event is a CAEP Event, defined here:
// https://openid.net/specs/openid-caep-specification-1_0-ID1.html#rfc.section.3.3
type CredentialChangeEvent struct {
	// Json defines the raw JSON of the CAEP Event. Used if
	// a developer wants greater control over all the attributes
	// of the CAEP Event
	Json map[string]interface{}

	// SubjectFormat defines the format of the subject event.
	//
	// See: https://datatracker.ietf.org/doc/html/draft-ietf-secevent-subject-identifiers
	// for an overview of available subject formats for CAEP Events
	Format SubjectFormat

	// Subject defines the subject that the CAEP Event applies to.
	//
	// See your transmitter's specification for the exact format
	// of the Subject
	Subject map[string]interface{}

	// EventTimestamp defines the timestamp of the CAEP Event in
	// Unix time (seconds since January 1, 1970 UTC)
	EventTimestamp int64

	// CredentialType defines the type of credential of the CAEP Event that has been modified/removed.
	// See https://openid.net/specs/openid-caep-specification-1_0.html#rfc.section.3.3.1 for the options for this field
	CredentialType CredentialType

	// CredentialType defines the type of modification/deletion towards the credential of the CAEP Event.
	// See https://openid.net/specs/openid-caep-specification-1_0.html#rfc.section.3.3.1 for the options for this field
	ChangeType ChangeType
}

func (event *CredentialChangeEvent) GetEventUri() string {
	return "https://schemas.openid.net/secevent/caep/event-type/credential-change"
}

func (event *CredentialChangeEvent) GetSubjectFormat() SubjectFormat {
	return event.Format
}

func (event *CredentialChangeEvent) GetSubject() map[string]interface{} {
	return event.Subject
}

func (event *CredentialChangeEvent) GetTimestamp() int64 {
	return event.EventTimestamp
}

func (event *CredentialChangeEvent) GetCredentialType() CredentialType {
	return event.CredentialType
}

func (event *CredentialChangeEvent) GetChangeType() ChangeType {
	return event.ChangeType
}

func (event *CredentialChangeEvent) GetType() EventType {
	return CredentialChange
}
