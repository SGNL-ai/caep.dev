package main

import (
	"fmt"
	"time"

	"github.com/sgnl-ai/caep.dev-receiver/pkg"
	ssf_events "github.com/sgnl-ai/caep.dev-receiver/pkg/ssf_events"
)

func main() {
	// Configure the receiver (specify the poll callback function and poll interval to start polling in each interval)
	receiverConfig := pkg.ReceiverConfig{
		TransmitterUrl:     "https://ssf.caep.dev",
		TransmitterPollUrl: "https://ssf.caep.dev/ssf/streams/poll",
		EventsRequested:    []ssf_events.EventType{0, 1, 2, 3, 4},
		AuthorizationToken: "<access token>",
		PollCallback:       PrintEvents,
		PollInterval:       20,
	}

	// Initialize the receiver and start polling
	receiver, err := pkg.ConfigureSsfReceiver(receiverConfig)
	if err != nil {
		print(err)
	}

	time.Sleep(time.Duration(90) * time.Second)
	// Delete the receiver after done
	receiver.DeleteReceiver()
}

func PrintEvents(events []ssf_events.SsfEvent) {
	fmt.Printf("Number of events: %d\n", len(events))
	for _, event := range events {
		fmt.Println("--------EVENT-------")
		fmt.Printf("Subject Format: %v\n", event.GetSubjectFormat())
		fmt.Printf("Subject: %v\n", event.GetSubject())
		fmt.Printf("Timestamp: %d\n", event.GetTimestamp())
		fmt.Printf("event.GetType(): %v\n", event.GetType())
		eventType := event.GetType()
		switch eventType {
		case ssf_events.AssuranceLevelChange:
			e, ok := event.(*ssf_events.AssuranceLevelChangeEvent)
			if !ok {
				fmt.Printf("cannot convert")
			}
			fmt.Printf("e: %v\n", e)

		case ssf_events.CredentialChange:
			e, ok := event.(*ssf_events.CredentialChangeEvent)
			if !ok {
				fmt.Printf("cannot convert")
			}
			fmt.Printf("e: %v\n", e)

		case ssf_events.DeviceCompliance:
			e, ok := event.(*ssf_events.DeviceComplianceEvent)
			if !ok {
				fmt.Printf("cannot convert")
			}
			fmt.Printf("e: %v\n", e)

		case ssf_events.SessionRevoked:
			e, ok := event.(*ssf_events.SessionRevokedEvent)
			if !ok {
				fmt.Printf("cannot convert")
			}
			fmt.Printf("e: %v\n", e)

		case ssf_events.TokenClaimsChange:
			e, ok := event.(*ssf_events.TokenClaimsChangeEvent)
			if !ok {
				fmt.Printf("cannot convert")
			}
			fmt.Printf("e: %v\n", e)

		case ssf_events.VerificationEventType:
			e, ok := event.(*ssf_events.VerificationEvent)
			if !ok {
				fmt.Printf("cannot convert")
			}
			fmt.Printf("e: %v\n", e)

		case ssf_events.StreamUpdatedEventType:
			e, ok := event.(*ssf_events.StreamUpdatedEvent)
			if !ok {
				fmt.Printf("cannot convert")
			}
			fmt.Printf("e: %v\n", e)
		}
		fmt.Println("--------------------")
	}
	fmt.Print("\n\n")
}
