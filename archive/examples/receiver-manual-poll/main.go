package main

import (
	"fmt"
	"time"

	"github.com/sgnl-ai/caep.dev-receiver/pkg"
	events "github.com/sgnl-ai/caep.dev-receiver/pkg/ssf_events"
)

func main() {
	// Configure the receiver (do not specify poll callback if polling is not required yet)
	receiverConfig := pkg.ReceiverConfig{
		TransmitterUrl:     "https://ssf.caep.dev",
		TransmitterPollUrl: "https://ssf.caep.dev/ssf/streams/poll",
		EventsRequested:    []events.EventType{0, 1, 2, 3, 4},
		AuthorizationToken: "<access token>",
		PollCallback:       nil,
	}
	// Initialize the receiver but does not start polling
	receiver, err := pkg.ConfigureSsfReceiver(receiverConfig)
	if err != nil {
		print(err)
	}

	print("sleep...\n")
	time.Sleep(time.Duration(10) * time.Second)
	print("awake!\n")
	// Manually call the PollEvents() to poll events at desired time
	events, err := receiver.PollEvents()
	if err != nil {
		print(err)
	}

	PrintEvents(events)
	// Delete the receiver after done
	receiver.DeleteReceiver()
	fmt.Println("Deleting receiver...")

}

func PrintEvents(events []events.SsfEvent) {
	fmt.Printf("Number of events: %d\n", len(events))
	for _, event := range events {
		fmt.Println("--------EVENT-------")
		fmt.Printf("Subject Format: %v\n", event.GetSubjectFormat())
		fmt.Printf("Subject: %v\n", event.GetSubject())
		fmt.Printf("Timestamp: %d\n", event.GetTimestamp())
		fmt.Println("--------------------")
	}
	fmt.Print("\n\n")
}
