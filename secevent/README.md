# secevent

A comprehensive Go library for building, signing, parsing, and validating Security Event Tokens (SecEvents) according to the Security Event Token [RFC 8417](https://tools.ietf.org/html/rfc8417).

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Building a SecEvent](#building-a-secevent)
  - [Parsing a SecEvent](#parsing-a-secevent)
- [Standard Events Support](#standard-events-support)
- [Defining Custom Events](#defining-custom-events)
- [Using Custom Signers](#using-custom-signers)
- [Subjects and Identifiers](#subjects-and-identifiers)
- [ID Generators](#id-generators)
- [Providing Verification Keys](#providing-verification-keys)
- [Contributing](#contributing)

---

## Features

- **Complete SecEvent Implementation**: Full support for building, signing, parsing, and validating Security Event Tokens in adherence to RFC 8417.
- **Out-of-the-box Support for Standard SecEvents**: Provides out-of-the-box support for CAEP and standard SSF events. Contributions for additional standard event support are welcome.
- **Event Extensibility**: Users can define event types for scenarios not covered by the library.
- **Flexible Subject Identifiers**: Supports various subject identifier formats, including email, phone number, issuer and subject pairs, URIs, and more.
- **Extensible Signing Mechanisms**: Integrate with custom signing functions or hardware security modules (HSMs) when private keys are not directly accessible.
- **Flexible Key Provisioning for Parsing**: Supports multiple ways to provide verification keys, such as JWKS URLs, JWKS JSON, or direct public keys.
- **Rich Examples**: Comprehensive examples demonstrating how to use the library for common and advanced use cases.

---

## Installation

```bash
go get github.com/sgnl-ai/caep.dev/secevent
```

---

## Quick Start

### Building a SecEvent

```go
package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "fmt"
    "time"

    "github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/caep"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/builder"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/signing"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/subject"
)

func main() {
    // Create a builder with configuration
    secEventBuilder := builder.NewBuilder(
        builder.WithDefaultIssuer("https://issuer.example.com"),
        builder.WithDefaultIDGenerator(id.NewUUIDGenerator()),
    )

    // Create a session revoked event
    sessionEvent := caep.NewSessionRevokedEvent().
        WithInitiatingEntity(caep.InitiatingEntityPolicy).
        WithReasonAdmin("en", "Security policy violation").
        WithEventTimestamp(time.Now().Unix())

    // Create a subject (e.g., email)
    userEmail := subject.NewEmailSubject("user@example.com")

    // Create a SecEvent using builder
    // Note: No need to specify issuer and ID as they come from builder defaults
    secEvent := secEventBuilder.NewSecEvent().
        WithAudience("https://receiver.example.com").
        WithSubject(userEmail).
        WithEvent(sessionEvent)

    // Generate a private key (for example purposes)
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

    fmt.Println("Signed SecEvent:", signedToken)

    // Create another SecEvent overriding the defaults
    nonDefaultsSecEvent := secEventBuilder.NewSecEvent().
        WithIssuer("https://custom-issuer.example.com").        // Override default issuer
        WithID("unique-sec-event-id").                         // Override generated ID
        WithAudience("https://receiver.example.com").
        WithSubject(userEmail).
        WithEvent(sessionEvent)
}
```

### Parsing a SecEvent

```go
package main

import (
    "fmt"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/parser"
)

func main() {
    tokenString := "..." // Your SecEvent JWT string

    // Create a parser with JWKS URL
    secEventParser := parser.NewParser(
        parser.WithJWKSURL("https://issuer.example.com/jwks.json"),
        parser.WithExpectedIssuer("https://issuer.example.com"),
        parser.WithExpectedAudience("https://receiver.example.com"),
    )

    // Parse SecEvent with signature verification
    secEvent, err := secEventParser.ParseSecEvent(tokenString)
    if err != nil {
        panic(err)
    }

    // Access SecEvent's event
    event := secEvent.Event
    fmt.Printf("Event Type: %s\n", event.Type())

    // Access SecEvent's subject
    subject := secEvent.Subject
    fmt.Printf("Subject Format: %s\n", subject.Format())

    // Parse without verification (useful for debugging)
    unverifiedSecEvent, err := secEventParser.ParseSecEventNoVerify(tokenString)
    if err != nil {
        panic(err)
    }
}
```

---

## Standard Events Support

The library provides supports for standard SecEvents. As of now, the following standard events are supported:

CAEP Events:
- `token-claims-change`
- `session-revoked`
- `credential-change`
- `assurance-level-change`
- `device-compliance-change`

SSF Events:
- `verification`
- `stream-updated`

**Example: Using Standard Events**

```go
package main

import (
    "github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/caep"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/ssf"
)

func main() {
    // Create CAEP event
    sessionRevokedEvent := caep.NewSessionRevokedEvent().
        WithInitiatingEntity(caep.InitiatingEntityPolicy).
        WithReasonAdmin("en", "Security policy violation")

    // Create SSF event
    verificationEvent := ssf.NewVerificationEvent().
        WithState("verification-state-123")

    // Create stream update event
    streamEvent := ssf.NewStreamUpdateEvent(ssf.StreamStatusEnabled).
        WithReason("Stream activated by admin")
}
```

---

## Defining Custom Events

The library allows you to define custom events.

**Example: Defining and Using a Custom Event**

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/sgnl-ai/caep.dev/secevent/pkg/event"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/builder"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/subject"
)

// Define your custom event type
const CustomEventType event.EventType = "https://example.com/event-type/custom"

// CustomEventPayload represents the payload for a custom event
type CustomEventPayload struct {
    CustomField string `json:"custom_field"`
}

// CustomEvent represents a custom event
type CustomEvent struct {
    event.BaseEvent
    CustomEventPayload
}

// NewCustomEvent creates a new custom event
func NewCustomEvent(customField string) *CustomEvent {
    e := &CustomEvent{
        CustomEventPayload: CustomEventPayload{
            CustomField: customField,
        },
    }

    e.SetType(CustomEventType)

    return e
}

// Validate ensures the event is valid
func (e *CustomEvent) Validate() error {
    if e.CustomField == "" {
        return event.NewError(event.ErrCodeMissingValue, "custom_field is required", "custom_field")
    }

    return nil
}

// Payload returns the event payload
func (e *CustomEvent) Payload() interface{} {
    return e.CustomEventPayload
}

// MarshalJSON implements the json.Marshaler interface
func (e *CustomEvent) MarshalJSON() ([]byte, error) {
    return json.Marshal(e.Payload())
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (e *CustomEvent) UnmarshalJSON(data []byte) error {
    var payload CustomEventPayload
    if err := json.Unmarshal(data, &payload); err != nil {
        return event.NewError(event.ErrCodeParseError, "failed to parse custom event data", "")
    }

    e.SetType(CustomEventType)

    e.CustomEventPayload = payload

    return e.Validate()
}

func parseCustomEvent(data []byte) (event.Event, error) {
    var e CustomEvent
    if err := json.Unmarshal(data, &e); err != nil {
        return nil, err
    }

    return &e, nil
}

func init() {
    // Register the event parser
    event.RegisterEventParser(CustomEventType, parseCustomEvent)
}


func main() {
    // Create a custom event
    customEvent := NewCustomEvent("Custom Value")

    // Create a subject
    userEmail := subject.NewEmailSubject("user@example.com")

    // Build the SecEvent
    secEvent := builder.NewSecEvent().
        WithIssuer("https://issuer.example.com").
        WithID("unique-id").
        WithSubject(userEmail).
        WithEvent(customEvent)

    // ... continue with signing if needed
}
```

---

## Using Custom Signers

The library supports custom signing mechanisms, allowing integration with HSMs or external signing services where private keys are not directly accessible.

**Implementing a Custom Signer**

```go
package main

import (
    "strings"
    "github.com/golang-jwt/jwt/v5"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/signing"
)

// CustomSigner implements the Signer interface
type CustomSigner struct {
    // Fields for HSM client or external service configuration
}

func (s *CustomSigner) Sign(claims jwt.Claims) (string, error) {
    // Create token with claims
    token := jwt.NewWithClaims(s.signingMethod, claims)
    
    // Set required headers for SecEvents
    token.Header["kid"] = s.kid  // Key ID
    token.Header["typ"] = "secevent+jwt"  // Token type for SecEvents
    
    // Obtain the signing string
    signingString, err := token.SigningString()
    if err != nil {
        return "", err
    }
    
    // Use your HSM or external service to sign the string
    signature, err := externalSign(signingString)
    if err != nil {
        return "", err
    }
    
    // Return the complete JWT
    return strings.Join([]string{signingString, signature}, "."), nil
}

// externalSign is a placeholder function representing the external signing process
func externalSign(signingString string) (string, error) {
    // Implement the signing logic using your HSM or external service
    return "signature", nil
}

func main() {
    customSigner := &CustomSigner{}

    // Use the custom signer with your SecEvent
    secEvent := builder.NewSecEvent().
        WithIssuer("https://issuer.example.com")
        // ... other SecEvent configuration

    signedToken, err := customSigner.Sign(secEvent)
    if err != nil {
        panic(err)
    }
}
```

## Subjects and Identifiers

The library supports various subject identifier formats as defined in RFC 9493.

**Supported Formats**

- `email`: Identified by an email address
- `phone_number`: Identified by a phone number
- `iss_sub`: Identified by an issuer and subject pair
- `uri`: Identified by a URI
- `opaque`: Identified by an opaque identifier
- `account`: Identified by an `acct:` URI
- `did`: Identified by a Decentralized Identifier (DID) URL
- `jwt_id`: Identified by a JWT issuer and ID
- `saml_assertion_id`: Identified by a SAML assertion issuer and ID
- `complex`: A composite subject made up of multiple components

**Example: Using Different Subject Types**

```go
package main

import (
    "github.com/sgnl-ai/caep.dev/secevent/pkg/subject"
)

func main() {
    // Email subject
    emailSubject := subject.NewEmailSubject("user@example.com")

    // Phone number subject
    phoneSubject := subject.NewPhoneSubject("+1-555-123-4567")

    // Issuer and subject pair
    issSubSubject := subject.NewIssSubSubject("https://issuer.example.com", "user123")

    // Complex subject combining multiple identifiers
    complexSubject := subject.NewComplexSubject().
        WithUser(emailSubject).
        WithDevice(subject.NewOpaqueSubject("device-123"))

    // Use these subjects when creating SecEvents
}
```

---

## ID Generators

Customize how SecEvent IDs (`jti`) are generated.

**Built-in Generators**

- `id.UUID`: Generates UUIDs
- `id.Random32`: Generates 32-character random hex strings
- `id.Base64URL20`: Generates 20-character base64url-encoded random strings
- `id.Sequential`: Generates sequential IDs with a prefix

**Example: Using Custom ID Generator**

```go
package main

import (
    "github.com/sgnl-ai/caep.dev/secevent/pkg/id"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/builder"
)

// Define a custom generator
type customGenerator struct {}

func (g *customGenerator) Generate() string {
    return "custom-prefix-" + time.Now().Format("20060102150405")
}

func main() {
    // Use built-in UUID generator
    secEvent1 := builder.NewSecEvent().
        WithIDGenerator(id.NewUUIDGenerator())

    // Use custom generator
    customGen := &customGenerator{}
    secEvent2 := builder.NewSecEvent().
        WithIDGenerator(customGen)
}
```

---

## Providing Verification Keys

When parsing and verifying SecEvents, you can provide verification keys in multiple ways:

**Using JWKS URL**

```go
parser := builder.NewParser().
    WithJWKSURL("https://issuer.example.com/jwks.json").
    WithExpectedIssuer("https://issuer.example.com")
```

**Using JWKS JSON**

```go
jwksJSON := []byte(`{"keys": [...]}`) // Replace with actual JWKS JSON
parser := builder.NewParser().
    WithJWKSJSON(jwksJSON).
    WithExpectedIssuer("https://issuer.example.com")
```

**Using Direct Public Key**

```go
publicKey := /* Load your public key */
parser := builder.NewParser().
    WithPublicKey(publicKey).
    WithExpectedIssuer("https://issuer.example.com")
```

---

## Contributing

Contributions to the project are welcome, including feature enhancements, bug fixes, and documentation improvements.