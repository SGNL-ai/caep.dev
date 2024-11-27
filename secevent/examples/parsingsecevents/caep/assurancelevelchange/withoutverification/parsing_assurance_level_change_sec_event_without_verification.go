package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/sgnl-ai/caep.dev/secevent/pkg/builder"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/id"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/parser"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/caep"
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
	signedSecEvent, err := generateSignedAssuranceLevelChangeSecEvent(keyPair.PrivateKey)
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
	if caepEvent, ok := event.(*caep.AssuranceLevelChangeEvent); ok {
		fmt.Printf("Current Level: %s\n", caepEvent.CurrentLevel)
		fmt.Printf("Previous Level: %s\n", caepEvent.PreviousLevel)
		fmt.Printf("Change Direction: %s\n", caepEvent.ChangeDirection)
		if initiatingEntity, exist := caepEvent.GetInitiatingEntity(); exist {
			fmt.Printf("Initiating Entity: %s\n", initiatingEntity)
		}

		if reasonAdmins := caepEvent.GetAllReasonAdmin(); reasonAdmins != nil {
			for lang, text := range reasonAdmins {
				fmt.Printf("Admin Reason (%s): %s\n", lang, text)
			}
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

func generateSignedAssuranceLevelChangeSecEvent(privateKey *ecdsa.PrivateKey) (string, error) {
	// Create a builder with configuration
	secEventBuilder := builder.NewBuilder(
		builder.WithDefaultIssuer("https://issuer.example.com"),
		builder.WithDefaultIDGenerator(id.NewUUIDGenerator()),
	)

	// Create an Assurance Level Change Event
	assuranceLevelChangeEvent := caep.NewAssuranceLevelChangeEvent(
		caep.AssuranceLevelAAL1,
		caep.AssuranceLevelAAL2,
		caep.ChangeDirectionDecrease,
	).
		WithInitiatingEntity(caep.InitiatingEntityPolicy).
		WithReasonAdmin("en", "Security policy violation").
		WithEventTimestamp(time.Now().Unix())

	// Create a subject (email)
	userEmail, err := subject.NewEmailSubject("user@example.com")
	if err != nil {
		return "", fmt.Errorf("failed to create subject: %w", err)
	}

	// Create a SecEvent using builder
	secEvent := secEventBuilder.NewSecEvent().
		WithAudience("https://receiver.example.com").
		WithSubject(userEmail).
		WithEvent(assuranceLevelChangeEvent)

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
