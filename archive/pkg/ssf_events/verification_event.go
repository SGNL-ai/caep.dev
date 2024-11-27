package ssf_events

// The Verification event is an SSE Event, defined here:
// https://openid.github.io/sharedsignals/openid-sharedsignals-framework-1_0.html#name-verification
type VerificationEvent struct {
	// Json defines the raw JSON of the CAEP Event. Used if
	// a developer wants greater control over all the attributes
	// of the CAEP Event
	Json map[string]interface{}

	// Event Receivers MAY use the value of this parameter to
	// correlate a verification event with a verification request.
	State string
}

func (event *VerificationEvent) GetEventUri() string {
	return "https://schemas.openid.net/secevent/ssf/event-type/verification"
}

func (event *VerificationEvent) GetSubjectFormat() SubjectFormat {
	return -1
}

func (event *VerificationEvent) GetSubject() map[string]interface{} {
	return map[string]interface{}{}
}

func (event *VerificationEvent) GetTimestamp() int64 {
	return 0
}

func (event *VerificationEvent) GetType() EventType {
	return VerificationEventType
}

func (event *VerificationEvent) GetState() string {
	return event.State
}
