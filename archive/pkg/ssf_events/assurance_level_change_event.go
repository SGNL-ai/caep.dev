package ssf_events

// The session revoked event is a CAEP Event, defined here:
// https://openid.net/specs/openid-caep-specification-1_0-ID1.html#rfc.section.3.1
type AssuranceLevelChangeEvent struct {
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

	// the namespace of the values in the current_level and previous_level claims.
	// See: https://openid.github.io/sharedsignals/openid-caep-specification-1_0.html#name-event-specific-claims-4 for valid values.
	Namespace string

	// the current NIST Authenticator Assurance Level (AAL) as defined in [SP800-63R3].
	// See: https://openid.github.io/sharedsignals/openid-caep-specification-1_0.html#name-event-specific-claims-4 for valid values.
	CurrentLevel string

	// the previous NIST Authenticator Assurance Level (AAL) as defined in [SP800-63R3]. Optional.
	// See: https://openid.github.io/sharedsignals/openid-caep-specification-1_0.html#name-event-specific-claims-4 for valid values.
	PreviousLevel *string

	// the Authenticator Assurance Level increased or decreased. Optional.
	// Must be either 'increase' or 'decrease'. See: See: https://openid.github.io/sharedsignals/openid-caep-specification-1_0.html#name-event-specific-claims-4
	// for more details.
	ChangeDirection *string
}

func (event *AssuranceLevelChangeEvent) GetEventUri() string {
	return "https://schemas.openid.net/secevent/caep/event-type/assurance-level-change"
}

func (event *AssuranceLevelChangeEvent) GetSubjectFormat() SubjectFormat {
	return event.Format
}

func (event *AssuranceLevelChangeEvent) GetSubject() map[string]interface{} {
	return event.Subject
}

func (event *AssuranceLevelChangeEvent) GetTimestamp() int64 {
	return event.EventTimestamp
}

func (event *AssuranceLevelChangeEvent) GetPreviousLevel() string {
	return *event.PreviousLevel
}

func (event *AssuranceLevelChangeEvent) GetCurrentLevel() string {
	return event.CurrentLevel
}

func (event *AssuranceLevelChangeEvent) GetChangeDirection() string {
	return *event.ChangeDirection
}

func (event *AssuranceLevelChangeEvent) GetNamespace() string {
	return event.Namespace
}

func (event *AssuranceLevelChangeEvent) GetType() EventType {
	return AssuranceLevelChange
}
