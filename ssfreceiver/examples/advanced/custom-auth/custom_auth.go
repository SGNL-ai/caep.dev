package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/parser"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/caep"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/ssf"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/builder"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/options"
)

// APIKeyAuth implements custom authorization
type APIKeyAuth struct {
	apiKey   string
	tenantID string
}

// AddAuth implements the auth.Authorizer interface
func (a *APIKeyAuth) AddAuth(ctx context.Context, req *http.Request) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Add custom headers for authentication
	req.Header.Set("X-API-Key", a.apiKey)
	req.Header.Set("X-Tenant-ID", a.tenantID)

	return nil
}

func main() {
	// Initialize custom authentication
	apiKeyAuth := &APIKeyAuth{
		apiKey:   "your-api-key-here",
		tenantID: "your-tenant-id",
	}

	streamBuilder, err := builder.New(
		"https://transmitter.example.com/.well-known/ssf-configuration", // transmitter's metadata url
		builder.WithPollDelivery(),
		builder.WithAuth(apiKeyAuth),
		builder.WithEventTypes([]event.EventType{
			caep.EventTypeSessionRevoked,
			caep.EventTypeDeviceComplianceChange,
		}),
		builder.WithExistingCheck(), // retrieve existing stream if exist
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

			case caep.EventTypeDeviceComplianceChange:
				log.Printf("Received device compliance event: %s", jti)

				if subjectPayload, err := secEvent.Subject.Payload(); err == nil {
					log.Printf("Subject details: %+v", subjectPayload)
				} else {
					log.Printf("Failed to get subject payload: %v", err)
				}

				if complianceEvent, ok := secEvent.Event.(*caep.DeviceComplianceChangeEvent); ok {
					log.Printf("Current Status: %s", complianceEvent.GetCurrentStatus())
					log.Printf("Previous Status: %s", complianceEvent.GetPreviousStatus())
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
