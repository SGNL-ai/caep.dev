package pkg

import event "github.com/sgnl-ai/caep.dev-receiver/pkg/ssf_events"

// Represents the interface for the SSF receiver with user facing
// methods
type SsfReceiver interface {
	ConfigureCallback(callback func(events []event.SsfEvent), pollInterval int) error

	// Polls the configured receiver a returns a list of the available SSF
	// Events
	PollEvents() ([]event.SsfEvent, error)

	// Cleans up the Receiver's resources and deletes it from the transmitter
	DeleteReceiver()

	// Get stream status from the transmitter
	GetStreamStatus() (StreamStatus, error)

	// Enable the stream
	EnableStream() (StreamStatus, error)

	// Pause the stream
	PauseStream() (StreamStatus, error)

	// Disable the stream
	DisableStream() (StreamStatus, error)
}

// The struct that contains all the necessary fields and methods for the
// SSF Receiver's implementation
type SsfReceiverImplementation struct {
	// transmitterUrl stores the base url of the transmitter the
	// receiver will make request to
	transmitterUrl string

	// transmitterPollUrl defines the url that the receiver
	// should hit to receive SSF Events
	transmitterPollUrl string

	// TransmitterStreamUrl defines the URL that the receiver will use
	// to update/get the stream status
	transmitterStatusUrl string

	// eventsRequested contains a list of the SSF Event URI's requested
	// by the receiver
	eventsRequested []string

	// authorizationToken defines the Auth Token used to authorize the
	// receiver with the transmitter
	authorizationToken string

	// pollCallback defines the method the receiver will call to pass
	// events into when the poll interval is triggered
	pollCallback func(events []event.SsfEvent)

	// pollInterval defines the interval, in seconds, between every
	// poll request the receiver will make to the transmitter. After
	// each poll request, the available SSF events will be passed in
	// a function call to pollCallback
	pollInterval int

	// configurationUrl defines the transmitter's configuration url
	configurationUrl string

	// streamId defines the Id of the stream that corresponds to the
	// transmitter
	streamId string

	// terminate is used to stop the push interval routine
	terminate chan bool
}

// Struct used to read a Transmitter's configuration
type TransmitterConfig struct {
	Issuer                   string                   `json:"issuer"`
	JwksUri                  string                   `json:"jwks_uri,omitempty"`
	DeliveryMethodsSupported []string                 `json:"delivery_methods_supported,omitempty"`
	ConfigurationEndpoint    string                   `json:"configuration_endpoint,omitempty"`
	StatusEndpoint           string                   `json:"status_endpoint,omitempty"`
	SpecVersion              string                   `json:"spec_version,omitempty"`
	AuthorizationSchemes     []map[string]interface{} `json:"authorization_schemes,omitempty"`
}

// Struct used to make a Create Stream request for the receiver
type CreateStreamReq struct {
	Delivery        SsfDelivery `json:"delivery"`
	EventsRequested []string    `json:"events_requested"`
	Description     string      `json:"description,omitempty"`
}

// Struct that defines the deliver method for the Create Stream Request
type SsfDelivery struct {
	Method string `json:"method"`
}

// Struct to make a request to poll SSF Events to the
// configured transmitter
type PollTransmitterRequest struct {
	Acknowledgements  []string `json:"ack"`
	MaxEvents         int      `json:"maxEvents,omitempty"`
	ReturnImmediately bool     `json:"returnImmediately"`
}

// Struct to make a request to update the stream status
type UpdateStreamRequest struct {
	StreamId string `json:"stream_id"`
	Status   string `json:"status"`
	Reason   string `json:"reason"`
}

type StreamStatus int

const (
	StreamEnabled StreamStatus = iota + 1
	StreamPaused
	StreamDisabled
)

var StatusEnumMap = map[string]StreamStatus{
	"enabled":  StreamEnabled,
	"paused":   StreamPaused,
	"disabled": StreamDisabled,
}

var EnumToStringStatusMap = map[StreamStatus]string{
	StreamEnabled:  "enabled",
	StreamPaused:   "paused",
	StreamDisabled: "disabled",
}
