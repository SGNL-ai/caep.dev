package ssf_events

// The session revoked event is a CAEP Event, defined here:
// https://openid.net/specs/openid-caep-specification-1_0-ID1.html#rfc.section.3.1
type SessionRevokedEvent struct {
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
}

func (event *SessionRevokedEvent) GetEventUri() string {
	return "https://schemas.openid.net/secevent/caep/event-type/session-revoked"
}

func (event *SessionRevokedEvent) GetSubjectFormat() SubjectFormat {
	return event.Format
}

func (event *SessionRevokedEvent) GetSubject() map[string]interface{} {
	return event.Subject
}

func (event *SessionRevokedEvent) GetTimestamp() int64 {
	return event.EventTimestamp
}

func (event *SessionRevokedEvent) GetType() EventType {
	return SessionRevoked
}
