package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/builder"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/id"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/caep"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/signing"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/subject"
)

func main() {
	// Create a builder with configuration
	secEventBuilder := builder.NewBuilder(
		builder.WithDefaultIssuer("https://issuer.example.com"),
		builder.WithDefaultIDGenerator(id.NewUUIDGenerator()),
	)

	// Create a Session Revoked Event
	sessionRevokedEvent := caep.NewSessionRevokedEvent().
		WithInitiatingEntity(caep.InitiatingEntityPolicy).
		WithReasonAdmin("en", "Security policy violation").
		WithEventTimestamp(time.Now().Unix())

	// Create a subject (e.g., email)
	userEmail, _ := subject.NewEmailSubject("user@example.com")

	// Create a SecEvent using builder
	secEvent := secEventBuilder.NewSecEvent().
		WithAudience("https://receiver.example.com").
		WithSubject(userEmail).
		WithEvent(sessionRevokedEvent)

	// Generate a private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	// Create a signer
	signer, err := signing.NewSigner(privateKey,
		signing.WithKeyID("key-1"))
	if err != nil {
		panic(err)
	}

	// Sign the SecEvent
	signedToken, err := signer.Sign(context.Background(), secEvent)
	if err != nil {
		panic(err)
	}

	fmt.Println("Signed Session Revoked SecEvent:", signedToken)
}
