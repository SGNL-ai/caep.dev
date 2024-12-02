package main

import (
	"context"
	"log"
	"time"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/parser"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/caep"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/ssf"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/auth"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/builder"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/options"
)

func main() {
	bearerAuth, err := auth.NewBearer("your-token-here")
	if err != nil {
		log.Fatalf("Failed to create bearer auth: %v", err)
	}

	streamBuilder, err := builder.New(
		"https://transmitter.example.com/.well-known/ssf-configuration", // transmitter's metadata url
		builder.WithPollDelivery(),
		builder.WithAuth(bearerAuth),
		builder.WithEventTypes([]event.EventType{
			caep.EventTypeSessionRevoked,
			caep.EventTypeCredentialChange,
		}),
		builder.WithExistingCheck(), // retrive existing stream if exist
	)
	if err != nil {
		log.Fatalf("Failed to create stream builder: %v", err)
	}

	log.Printf("Setting up stream connection...")

	stream, err := streamBuilder.Setup(context.Background())
	if err != nil {
		log.Fatalf("Failed to setup stream: %v", err)
	}

	log.Printf("Stream setup completed successfully")

	// Initialize SEC event parser
	secEventParser := parser.NewParser()

	for {
		log.Printf("Polling for new events...")

		rawEventTokens, err := stream.Poll(context.Background(),
			options.WithMaxEvents(10),
			options.WithAutoAck(true),
		)
		if err != nil {
			log.Fatalf("Failed to poll events: %v", err)
		}

		for jti, rawEventToken := range rawEventTokens {
			log.Printf("Processing event with JTI: %s", jti)
			// Note: ParseSecEventNoVerify do not verify the signature and claims of the SecEvent.
			// In production setting, use ParseSecEvent
			secEvent, err := secEventParser.ParseSecEventNoVerify(rawEventToken)
			if err != nil {
				log.Printf("Failed to parse event %s: %v", jti, err)

				continue
			}

			switch secEvent.Event.Type() {
			case caep.EventTypeSessionRevoked:
				log.Printf("Received session revoked event: %s", jti)

				if subjectPayload, err := secEvent.Subject.Payload(); err == nil {
					log.Printf("Subject details: %+v", subjectPayload)
				} else {
					log.Printf("Failed to get subject payload: %v", err)
				}
			case caep.EventTypeCredentialChange:
				log.Printf("Received credential change event: %s", jti)

				if subjectPayload, err := secEvent.Subject.Payload(); err == nil {
					log.Printf("Subject details: %+v", subjectPayload)
				} else {
					log.Printf("Failed to get subject payload: %v", err)
				}

				if credEvent, ok := secEvent.Event.(*caep.CredentialChangeEvent); ok {
					log.Printf("Credential type: %s", credEvent.GetCredentialType())
					log.Printf("Change type: %s", credEvent.GetChangeType())
				}
			case ssf.EventTypeVerification:
				log.Printf("Received verification event: %s", jti)

				if subjectPayload, err := secEvent.Subject.Payload(); err == nil {
					log.Printf("Subject details: %+v", subjectPayload)
				} else {
					log.Printf("Failed to get subject payload: %v", err)
				}

				if verificationEvent, ok := secEvent.Event.(*ssf.VerificationEvent); ok {
					if state, exist := verificationEvent.GetState(); exist {
						log.Printf("State: %s", state)
					}
				}
			default:
				log.Printf("Received unknown event type: %s", secEvent.Event.Type())

				if subjectPayload, err := secEvent.Subject.Payload(); err == nil {
					log.Printf("Subject details: %+v", subjectPayload)
				} else {
					log.Printf("Failed to get subject payload: %v", err)
				}
			}
		}

		time.Sleep(time.Second * 3) // Poll interval
	}
}
