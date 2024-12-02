package event

import (
	"encoding/json"
	"fmt"
)

type EventType string

// Event is the interface that all event types must implement
type Event interface {
	// Type returns the event type URI
	Type() EventType

	// Validate checks if the event is valid
	Validate() error

	// Payload returns the event payload
	Payload() interface{}

	// MarshalJSON implements the json.Marshaler interface
	json.Marshaler

	// UnmarshalJSON implements the json.Unmarshaler interface
	json.Unmarshaler
}

type BaseEvent struct {
	eventType EventType
	payload   interface{}
}

func (e *BaseEvent) Type() EventType {
	return e.eventType
}

func (e *BaseEvent) SetType(eventType EventType) {
	e.eventType = eventType
}

func (e *BaseEvent) Payload() interface{} {
	return e.payload
}

func (e *BaseEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.payload)
}

func (e *BaseEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.payload)
}

// EventParserRegistry is a map of event types to their respective parsers
var EventParserRegistry = map[EventType]func([]byte) (Event, error){}

// RegisterEventParser registers an event type and its parser
func RegisterEventParser(eventType EventType, parser func([]byte) (Event, error)) {
	EventParserRegistry[eventType] = parser
}

// GetEventParser returns the parser for a given event type
func GetEventParser(eventType EventType) (func([]byte) (Event, error), bool) {
	parser, ok := EventParserRegistry[eventType]

	return parser, ok
}

// ParseEvent parses event data based on the registered event type
func ParseEvent(eventType EventType, data []byte) (Event, error) {
	parser, ok := GetEventParser(eventType)
	if !ok {
		return nil, NewError(
			ErrCodeInvalidEventType,
			fmt.Sprintf("no parser registered for event type: %s", eventType),
			"event_type",
			"",
		)
	}

	// Validate that the data is valid JSON
	if !json.Valid(data) {
		return nil, NewError(
			ErrCodeParseError,
			"invalid JSON data",
			"",
			"",
		)
	}

	// Use the registered parser to parse the event
	event, err := parser(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event: %w", err)
	}

	// Verify the event type matches
	if event.Type() != eventType {
		return nil, NewError(
			ErrCodeInvalidEventType,
			fmt.Sprintf("parsed event type %s does not match expected type %s",
				event.Type(), eventType),
			"event_type",
			"",
		)
	}

	return event, nil
}

// GetRegisteredEventTypes returns all registered event types
func GetRegisteredEventTypes() []EventType {
	types := make([]EventType, 0, len(EventParserRegistry))
	for eventType := range EventParserRegistry {
		types = append(types, eventType)
	}

	return types
}

// IsEventTypeRegistered checks if a parser is registered for the given event type
func IsEventTypeRegistered(eventType EventType) bool {
	_, ok := EventParserRegistry[eventType]

	return ok
}
