package main

import (
	"context"
	"log"
	"time"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/parser"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/caep"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/ssf"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/subject"
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
		"https://transmitter.example.com/.well-known/ssf-configuration",
		builder.WithPollDelivery(),
		builder.WithAuth(bearerAuth),
		builder.WithEventTypes([]event.EventType{
			caep.EventTypeSessionRevoked,
		}),
		builder.WithExistingCheck(),
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

	// Create and add subjects to stream
	// 1. Email subject
	emailSubject, err := subject.NewEmailSubject("user@example.com")
	if err != nil {
		log.Fatalf("Failed to create email subject: %v", err)
	}

	// 2. Phone subject
	phoneSubject, err := subject.NewPhoneSubject("+1-555-123-4567")
	if err != nil {
		log.Fatalf("Failed to create phone subject: %v", err)
	}

	// 3. Complex subject combining user and device
	deviceSubject, err := subject.NewOpaqueSubject("device-123")
	if err != nil {
		log.Fatalf("Failed to create device subject: %v", err)
	}

	complexSubject := subject.NewComplexSubject().
		WithUser(emailSubject).
		WithDevice(deviceSubject)

	// Add subjects to stream
	ctx := context.Background()

	log.Printf("Adding email subject...")
	if err := stream.AddSubject(ctx, emailSubject); err != nil {
		log.Printf("Failed to add email subject: %v", err)
	}

	log.Printf("Adding phone subject...")
	if err := stream.AddSubject(ctx, phoneSubject); err != nil {
		log.Printf("Failed to add phone subject: %v", err)
	}

	log.Printf("Adding complex subject...")
	if err := stream.AddSubject(ctx, complexSubject); err != nil {
		log.Printf("Failed to add complex subject: %v", err)
	}

	// Initialize SEC event parser
	secEventParser := parser.NewParser()

	for {
		log.Printf("Polling for new events...")

		rawEventTokens, err := stream.Poll(ctx,
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

				// Remove the subject if session is revoked
				if err := stream.RemoveSubject(ctx, secEvent.Subject); err != nil {
					log.Printf("Failed to remove subject: %v", err)
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
			}
		}

		// Cleanup after processing events
		// Example: Remove phone subject after certain condition
		if shouldRemovePhoneSubject() {
			log.Printf("Removing phone subject...")

			if err := stream.RemoveSubject(ctx, phoneSubject); err != nil {
				log.Printf("Failed to remove phone subject: %v", err)
			}
		}

		time.Sleep(time.Second * 3) // Poll interval
	}
}

// shouldRemovePhoneSubject is a placeholder for your business logic
func shouldRemovePhoneSubject() bool {
	// Implement your logic here
	return false
}
