package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/builder"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/id"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/parser"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/ssf"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/signing"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/subject"
)

type KeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	PublicPEM  []byte
}

func main() {
	// Generate key pair
	keyPair, err := generateKeyPair()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate key pair: %v", err))
	}

	// Generate signed event
	signedSecEvent, err := generateSignedStreamUpdateSecEvent(keyPair.PrivateKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate signed event: %v", err))
	}

	// Create a parser
	secEventParser := parser.NewParser(
		parser.WithExpectedIssuer("https://issuer.example.com"),
		parser.WithExpectedAudience("https://receiver.example.com"),
	)

	// Parse without verification for demonstration
	secEvent, err := secEventParser.ParseSecEventNoVerify(signedSecEvent)
	if err != nil {
		panic(fmt.Errorf("failed to parse event: %w", err))
	}

	// Display all fields
	fmt.Println("\nParsed SecEvent Details:")
	fmt.Printf("Issuer: %s\n", secEvent.Issuer)
	fmt.Printf("Audience: %v\n", secEvent.Audience)
	fmt.Printf("Issued At: %v\n", secEvent.IssuedAt)
	fmt.Printf("JWT ID: %s\n", secEvent.ID)

	// Display subject information
	subject := secEvent.Subject
	fmt.Printf("\nSubject Information:\n")
	fmt.Printf("Format: %s\n", subject.Format())

	jsonSubject, err := subject.MarshalJSON()
	if err != nil {
		panic(fmt.Errorf("error getting json subject: %w", err))
	}

	fmt.Printf("Json Subject: %s\n", jsonSubject)

	// Display event information
	event := secEvent.Event
	fmt.Printf("\nEvent Information:\n")
	fmt.Printf("Type: %s\n", event.Type())

	// Type assert to access CAEP-specific fields
	if caepEvent, ok := event.(*ssf.StreamUpdateEvent); ok {
		fmt.Printf("Stream Status: %s\n", caepEvent.GetStatus())

		if reason, exist := caepEvent.GetReason(); exist {
			fmt.Printf("Stream Update Reason: %s\n", reason)
		}
	}
}

func generateKeyPair() (*KeyPair, error) {
	// Generate ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Extract public key
	publicKey := &privateKey.PublicKey

	// Encode public key to PEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		PublicPEM:  publicPEM,
	}, nil
}

func generateSignedStreamUpdateSecEvent(privateKey *ecdsa.PrivateKey) (string, error) {
	// Create a builder with configuration
	secEventBuilder := builder.NewBuilder(
		builder.WithDefaultIssuer("https://issuer.example.com"),
		builder.WithDefaultIDGenerator(id.NewUUIDGenerator()),
	)

	// Create an Assurance Level Change Event
	streamUpdateEvent := ssf.NewStreamUpdateEvent(ssf.StreamStatusDisabled).
		WithReason("example-stream-disable-reason")

	// Create a opaque subject
	streamSubject, _ := subject.NewOpaqueSubject("example-stream-id")

	// Create a SecEvent using builder
	secEvent := secEventBuilder.NewSecEvent().
		WithAudience("https://receiver.example.com").
		WithSubject(streamSubject).
		WithEvent(streamUpdateEvent)

	// Create a signer
	signer, err := signing.NewSigner(privateKey,
		signing.WithKeyID("key-1"))
	if err != nil {
		return "", fmt.Errorf("failed to create signer: %w", err)
	}

	// Sign the SecEvent
	signedToken, err := signer.Sign(secEvent)
	if err != nil {
		return "", fmt.Errorf("failed to sign event: %w", err)
	}

	return signedToken, nil
}
