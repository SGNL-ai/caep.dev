package pkg

import events "github.com/sgnl-ai/caep.dev-receiver/pkg/ssf_events"

type ReceiverConfig struct {
	// TransmitterUrl defines the URL for the transmitter that
	// the configured receiver will create a stream with and receive
	// events from.
	//
	// Required
	TransmitterUrl string

	// TransmitterPollUrl defines the URL that the receiver will use
	// to poll for SSF events.
	//
	// Note - Must be a subpath of TransmitterUrl
	//
	// Required
	TransmitterPollUrl string

	// TransmitterStreamUrl defines the URL that the receiver will use
	// to update/get the stream status.
	//
	// Note - Must be a subpath of TransmitterUrl
	//
	// Optional
	TransmitterStreamUrl string

	// EventsRequested specified the SSF events you want to receiver
	// from the transmitter.
	//
	// Required
	EventsRequested []events.EventType

	// AuthorizationToken is the authorization token used to authorize
	// your receiver with the specified transmitter
	//
	// Note - all transmitter's will require an authorization token
	//
	// Required
	AuthorizationToken string

	// PollCallback is used to configure the method that you want the
	// receiver to call after each automatic poll request. Each time
	// the poll interval timer is up, the receiver will make a request
	// to the specified transmitter and fetch available SSF events. It
	// will then call PollCallback with a list of those events
	//
	// Note - The PollCallback and PollInterval can also be configured
	// after initial receiver construction
	//
	// Optional
	PollCallback func(events []events.SsfEvent)

	// PollInterval defines, in seconds how often you want the receiver to
	// poll for SSF events any and pass them to your PollCallback function.
	//
	// Note - This field will not be used if the PollCallback isn't configured
	//
	// Optional, defaults to 300 (5 minutes)
	PollInterval int
}
