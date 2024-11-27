package main

import (
	"fmt"

	"github.com/sgnl-ai/caep.dev-receiver/pkg"
	events "github.com/sgnl-ai/caep.dev-receiver/pkg/ssf_events"
)

func main() {
	// Configure the receiver.
	receiverConfig := pkg.ReceiverConfig{
		TransmitterUrl:     "https://ssf.caep.dev",
		TransmitterPollUrl: "https://ssf.caep.dev/ssf/streams/poll",
		EventsRequested:    []events.EventType{0, 1, 2, 3, 4},
		AuthorizationToken: "<access token>",
		PollCallback:       nil,
	}

	// Initialize the receiver but does not start polling.
	receiver, err := pkg.ConfigureSsfReceiver(receiverConfig)
	if err != nil {
		print(err)
	}

	// Get the receiver's stream status.
	streamStatus, err := receiver.GetStreamStatus()
	if err != nil {
		print(err)
	}
	fmt.Printf("streamStatus: %v\n", streamStatus)

	// Update the receiver's stream status.
	updatedStreamStatus, err := receiver.DisableStream()
	if err != nil {
		print(err)
	}
	fmt.Printf("updatedStreamStatus: %v\n", updatedStreamStatus)

	// Verify the receiver's stream status.
	streamStatus, err = receiver.GetStreamStatus()
	if err != nil {
		print(err)
	}
	fmt.Printf("streamStatus: %v\n", streamStatus)

	// Update the receiver's stream status back to enable.
	updatedStreamStatus, err = receiver.EnableStream()
	if err != nil {
		print(err)
	}
	fmt.Printf("updatedStreamStatus: %v\n", updatedStreamStatus)

	// Delete the receiver after done.
	receiver.DeleteReceiver()
	fmt.Println("Deleting receiver...")
}
