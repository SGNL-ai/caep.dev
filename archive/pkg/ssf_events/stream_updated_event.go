package ssf_events

// The Verification event is an SSE Event, defined here:
// https://openid.github.io/sharedsignals/openid-sharedsignals-framework-1_0.html#name-stream-updated-event
type StreamUpdatedEvent struct {
	// Json defines the raw JSON of the CAEP Event. Used if
	// a developer wants greater control over all the attributes
	// of the CAEP Event
	Json map[string]interface{}

	// Status defines the new status of the stream.
	Status string

	// Provides a short description of why the Transmitter has updated the status.
	Reason string
}

func (event *StreamUpdatedEvent) GetEventUri() string {
	return "https://schemas.openid.net/secevent/caep/event-type/stream-updated"
}

func (event *StreamUpdatedEvent) GetSubjectFormat() SubjectFormat {
	return -1
}

func (event *StreamUpdatedEvent) GetSubject() map[string]interface{} {
	return map[string]interface{}{}
}

func (event *StreamUpdatedEvent) GetTimestamp() int64 {
	return 0
}

func (event *StreamUpdatedEvent) GetType() EventType {
	return StreamUpdatedEventType
}

func (event *StreamUpdatedEvent) GetState() string {
	return event.Status
}
