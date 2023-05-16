package pkg

type SessionRevokedEvent struct {
	// Json defines the raw JSON of the CAEP Event. Used if
	// a developer wants greater control over all the attributes
	// of the CAEP Event
	Json map[string]interface{}

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

func (event *SessionRevokedEvent) GetSubject() map[string]interface{} {
	return event.Subject
}

func (event *SessionRevokedEvent) GetTimestamp() int64 {
	return event.EventTimestamp
}
