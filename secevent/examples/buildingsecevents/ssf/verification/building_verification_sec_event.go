package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/builder"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/id"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/ssf"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/signing"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/subject"
)

func main() {
	// Create a builder with configuration
	secEventBuilder := builder.NewBuilder(
		builder.WithDefaultIssuer("https://issuer.example.com"),
		builder.WithDefaultIDGenerator(id.NewUUIDGenerator()),
	)

	// Create a Assurance Level Change Event
	verificationEvent := ssf.NewVerificationEvent().
		WithState("example-state")

	// Create a subject (e.g., email)
	streamSubject, _ := subject.NewOpaqueSubject("example-stream-id")

	// Create a SecEvent using builder
	secEvent := secEventBuilder.NewSecEvent().
		WithAudience("https://receiver.example.com").
		WithSubject(streamSubject).
		WithEvent(verificationEvent)

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
	signedToken, err := signer.Sign(secEvent)
	if err != nil {
		panic(err)
	}

	fmt.Println("Signed Verification SecEvent:", signedToken)
}
