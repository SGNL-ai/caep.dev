package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/event"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/token"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/parser"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/caep"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/ssf"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/auth"
	"github.com/sgnl-ai/caep.dev/ssfreceiver/builder"
)

func main() {
	bearerAuth, err := auth.NewBearer("your-token-here")
	if err != nil {
		log.Fatalf("Failed to create bearer auth: %v", err)
	}

	streamBuilder, err := builder.New(
		"https://transmitter.example.com/.well-known/ssf-configuration", 	// transmitter's metadata url
		builder.WithPushDelivery("https://receiver.example.com/events"), 	// event-delivery endpoint
		builder.WithAuth(bearerAuth),
		builder.WithEventTypes([]event.EventType{
			caep.EventTypeSessionRevoked,
			caep.EventTypeAssuranceLevelChange,
		}),
		builder.WithExistingCheck(),	// retrive existing stream if exist
	)
	if err != nil {
		log.Fatalf("Failed to create stream builder: %v", err)
	}

	log.Printf("Setting up stream connection...")

	stream, err := streamBuilder.Setup(context.Background())
	if err != nil {
		log.Fatalf("Failed to setup stream: %v", err)
	}

	log.Printf("Stream setup completed successfully. Stream ID: %s", stream.GetStreamID())
}

// HandlePushedEvent processes incoming events.
// Use this in your event-delivery endpoint handler:
//
// http.HandleFunc("/events", HandlePushedEvent)
func HandlePushedEvent(w http.ResponseWriter, r *http.Request) {
	secEventParser := parser.NewParser()

	// Read the request rawEventToken
	rawEventToken, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	defer r.Body.Close()

	// Note: ParseSecEventNoVerify do not verify the signature and claims of the SecEvent.
	// In production setting, use ParseSecEvent
	secEvent, err := secEventParser.ParseSecEventNoVerify(string(rawEventToken))
	if err != nil {
		log.Printf("Failed to parse event: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	switch secEvent.Event.Type() {
	case caep.EventTypeSessionRevoked:
		handleSessionRevoked(secEvent)
	case caep.EventTypeAssuranceLevelChange:
		handleAssuranceLevelChange(secEvent)
	case ssf.EventTypeVerification:
		handleVerification(secEvent)
	default:
		log.Printf("Received unknown event type: %s", secEvent.Event.Type())
	}

	w.WriteHeader(http.StatusOK)
}

func handleSessionRevoked(secEvent *token.SecEvent) {
	log.Printf("Received session revoked event: %s", secEvent.ID)
				
	if subjectPayload, err := secEvent.Subject.Payload(); err == nil {
		log.Printf("Subject details: %+v", subjectPayload)
	} else {
		log.Printf("Failed to get subject payload: %v", err)
	}
}

func handleAssuranceLevelChange(secEvent *token.SecEvent) {
	log.Printf("Received assurance level change event: %s", secEvent.ID)

	if subjectPayload, err := secEvent.Subject.Payload(); err == nil {
		log.Printf("Subject details: %+v", subjectPayload)
	} else {
		log.Printf("Failed to get subject payload: %v", err)
	}

	if assuranceEvent, ok := secEvent.Event.(*caep.AssuranceLevelChangeEvent); ok {
		log.Printf("Current Level: %s", assuranceEvent.GetCurrentLevel())
		log.Printf("Previous Level: %s", assuranceEvent.GetPreviousLevel())
		log.Printf("Change Direction: %s", assuranceEvent.GetChangeDirection())
	}
}

func handleVerification(secEvent *token.SecEvent) {
	log.Printf("Received verification event: %s", secEvent.ID)

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
}
