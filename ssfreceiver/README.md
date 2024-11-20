# ssfreceiver

A Go library for implementing [Shared Signals Framework (SSF)](https://openid.github.io/sharedsignals/openid-sharedsignals-framework-1_0.html) receivers.

## Table of Contents

- [Features](#features)
- [Dependencies](#dependencies)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Authorization Setup](#authorization-setup)
  - [Poll-Based Stream](#poll-based-stream)
  - [Push-Based Stream](#push-based-stream)
- [Stream Management](#stream-management)
  - [Stream Creation](#stream-creation)
  - [Stream Configuration](#stream-configuration)
  - [Stream Status](#stream-status)
  - [Stream Verification](#stream-verification)
- [Event Handling](#event-handling)
  - [Polling Events](#polling-events)
  - [Events Acknowledgment](#events-acknowledgment)
- [Subject Management](#subject-management)
- [Authorization](#authorization)
- [Custom Events](#custom-events)
- [Best Practices](#best-practices)
- [Contributing](#contributing)

## Features

- **Complete SSF Receiver Implementation**: Full support for creating and managing SSF streams according to the OpenID SSF specification
- **Integrates with secevent**: Integrates with the [secevent](https://github.com/SGNL-ai/caep.dev/tree/AddGoSsfSetLibrary/secevent) library for SET handling, subject management, and event parsing
- **Multiple Delivery Methods**: Support for both poll-based and push-based event delivery
- **Flexible Authorization**: Built-in support for various authorization methods with per-operation overrides
- **Stream Lifecycle Management**: Comprehensive stream status control (enable, disable, pause)
- **Stream Verification**: Built-in support for stream verification

## Dependencies

This library requires:
- [secevent](https://github.com/SGNL-ai/caep.dev/tree/AddGoSsfSetLibrary/secevent)

## Installation

```bash
go get github.com/cape.dev-receiver/ssfreceiver
```

## Quick Start

### Authorization Setup

First, set up your authorization method:

```go
package main

import (
    "github.com/cape.dev-receiver/ssfreceiver/auth"
)

// Bearer token authentication
bearerAuth := auth.NewBearer("token")

// OAuth2 client credentials
oauth2Auth := auth.NewOAuth2ClientCredentials(auth.OAuth2Config{
    ClientID:     "client_id",
    ClientSecret: "client_secret",
    TokenURL:     "token_url",
    Scopes:       []string{"scope1", "scope2"},
})

// Custom authorization implementation
type CustomAuth struct {
    // Custom auth fields
}

func NewCustomAuth(...) *CustomAuth {
    // Initialize and return custom auth implementation
}

func (a *CustomAuth) AddAuth(req *http.Request) error {
    // custom auth logic

    req.Header.Set(header, value)

    return nil
}

// Initialize custom authentication with required parameters
customAuth := NewCustomAuth(...)
```

### Poll-Based Stream

```go
package main

import (
    "context"
    "log"
    "github.com/cape.dev-receiver/ssfreceiver/builder"
    "github.com/cape.dev-receiver/ssfreceiver/options"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/subject"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/caep"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

func main() {
    // Create a stream builder
    streamBuilder := builder.New("https://transmitter.example.com/.well-known/ssf-configuration",
        builder.WithPollDelivery(),
        builder.WithAuth(bearerAuth),
        builder.WithEventTypes([]event.EventType{
            caep.EventTypeSessionRevoked,      
            caep.EventTypeTokenClaimsChange,
            event.EventType("https://custom.example.com/events/custom"), // Custom event type
        }),
    )

    // Build the stream
    ssfStream, err := streamBuilder.Build(context.Background())
    if err != nil {
        // Handle error
    }

    defer ssfStream.Disable(context.Background())

    emailSubject, err := subject.NewEmailSubject("user@example.com")
    if err != nil {
        // Handle error
    }

    err = ssfStream.AddSubject(context.Background(), 
        emailSubject,
        options.WithSubjectVerification(true),
        options.WithAuth(customAuth), // Operation-specific auth override
    )
    if err != nil {
        // Handle error
    }

    // Poll for events
    events, err := ssfStream.Poll(context.Background(),
        options.WithMaxEvents(10),          // Default is 100
        options.WithAutoAck(true),          // Auto acknowledge received events
        options.WithAckJTIs([]string{"jti-1", "jti-2"}), // Acknowledge previous events
    )

    // Process events using secevent parser
    secEventParser := parser.NewParser(
        parser.WithJWKSURL("https://issuer.example.com/jwks.json"),
        parser.WithExpectedIssuer("https://issuer.example.com"),
        parser.WithExpectedAudience("https://receiver.example.com"),
    )

    for _, rawEvent := range events {
        parsedEvent, err := secEventParser.ParseSingleEventSecEvent(rawEvent)
        if err != nil {
            // Handle error
            
            continue
        }
        
        handleEvent(parsedEvent)
    }
}

func handleEvent(event *secevent.Event) {
    switch event.Type() {
    case caep.EventTypeSessionRevoked:
        handleSessionRevoked(event)
    case caep.EventTypeTokenClaimsChange:
        handleTokenClaimsChange(event)
    default:
        handleUnknownEvent(event)
    }
}
```

### Push-Based Stream

```go
package main

import (
    "context"
    "net/http"
    "github.com/cape.dev-receiver/ssfreceiver/builder"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/parser"
    "github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/caep"
)

func main() {
    // Create a stream builder
    streamBuilder := builder.New("https://transmitter.example.com/.well-known/ssf-configuration",
        builder.WithPushDelivery("https://receiver.example.com/events"),
        builder.WithAuth(bearerAuth),
        builder.WithEventTypes([]event.EventType{
            caep.EventTypeSessionRevoked,      
            caep.EventTypeTokenClaimsChange,
            event.EventType("https://custom.example.com/events/custom"), // Custom event type
        }),
    )

    // Build and validate the stream
    ssfStream, err := streamBuilder.Build(context.Background())
    if err != nil {
        // Handle error
    }

    defer ssfStream.Disable(context.Background())

    // For push-based delivery, use secevent parser to handle incoming events
    secEventParser := parser.NewParser(
        parser.WithJWKSURL("https://issuer.example.com/jwks.json"),
        parser.WithExpectedIssuer("https://issuer.example.com"),
        parser.WithExpectedAudience("https://receiver.example.com"),
    )

    // Set up HTTP handler for events
    http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
        event, err := secEventParser.ParseSingleEventSecEvent(r.Body)
        if err != nil {
            // Handle error
        }
        
        handleEvent(event)

        w.WriteHeader(http.StatusOK)
    })

    http.ListenAndServe(":8888", nil)
}
```

## Stream Management

### Stream Creation
```go
// Create builder with options
builder := builder.New(transmitterURL,
    builder.WithPollDelivery(),  // or WithPushDelivery(endpoint)
    builder.WithAuth(auth.NewBearer("default-token")),
    builder.WithEventTypes([]event.EventType{
        caep.EventTypeSessionRevoked,
        caep.EventTypeTokenClaimsChange,
    }),
    builder.WithRetryConfig(options.RetryConfig{
        MaxRetries: 3,
        BackoffDuration: time.Second * 2,
    }),
)

// Build stream
stream, err := builder.Build(ctx)
```

### Stream Configuration
```go
// Get current configuration
config, err := stream.GetConfig(ctx)

// Update stream configuration
err := stream.UpdateConfig(ctx,
    options.WithEventTypes([]event.EventType{
        caep.EventTypeSessionRevoked,
        caep.EventTypeTokenClaimsChange,
    }),
    options.WithDescription("Updated configuration"))
```

### Stream Status
```go
// Pause stream - temporarily stop receiving events
err := stream.Pause(ctx, 
    options.WithStatusReason("System maintenance"),
    options.WithAuth(customAuth))

// Disable stream - stop receiving events and clear queue
err := stream.Disable(ctx)

// Delete stream - completely remove the stream
err := stream.Delete(ctx)

// Enable stream - start or resume receiving events
err := stream.Enable(ctx)
```

### Stream Verification
```go
// Verify stream with state
err := stream.Verify(ctx,
    options.WithState("verification-state"),
    options.WithAuth(customAuth))
```

## Event Handling

### Polling Events
```go
// Poll with various options
events, err := stream.Poll(ctx,
    options.WithMaxEvents(10),          // Default is 100
    options.WithAutoAck(true),          // Auto acknowledge received events
    options.WithAckJTIs([]string{...}), // Acknowledge previous events
    options.WithAuth(customAuth))       // Override default auth

// Process events using secevent parser
secEventParser := secevent.NewParser(
    parser.WithJWKSURL("https://issuer.example.com/jwks.json"),
    parser.WithJWKSURL("https://issuer.example.com/jwks.json"),
    parser.WithExpectedIssuer("https://issuer.example.com"),
)

for _, rawEvent := range events {
    parsedEvent, err := secEventParser.ParseSingleEventSecEvent(rawEvent)
    if err != nil {
        // Handle error
        continue
    }
    
    // Handle the event
    switch parsedEvent.Event.Type() {
    case caep.EventTypeSessionRevoked:
        handleSessionRevoked(parsedEvent)
    case caep.EventTypeTokenClaimsChange:
        handleTokenClaimsChange(parsedEvent)
    default:
        handleUnknownEvent(parsedEvent)
    }
}
```

### Events Acknowledgment
```go
// Acknowledge specific events
err := stream.Acknowledge(ctx, 
    []string{"jti1", "jti2"},
    options.WithAuth(customAuth))
```

## Subject Management

Using secevent's subject package for subject creation and management:

```go
import "github.com/sgnl-ai/caep.dev/secevent/pkg/subject"

// Create various subject types using secevent
emailSubject, err := subject.NewEmailSubject("user@example.com")
phoneSubject, err := subject.NewPhoneSubject("+1-555-123-4567")
issSubSubject, err := subject.NewIssSubSubject("https://issuer.example.com", "user123")

// Complex subject
complexSubject := subject.NewComplexSubject().
    WithUser(emailSubject).
    WithDevice(subject.NewOpaqueSubject("device-123"))

// Add subject to stream
err := stream.AddSubject(ctx,
    emailSubject,
    options.WithSubjectVerification(true),
    options.WithAuth(customAuth))

// Remove subject from stream
err := stream.RemoveSubject(ctx,
    emailSubject,
    options.WithAuth(customAuth))
```

## Authorization

The library provides a flexible authorization system through the `auth` package:

```go
package auth

// Authorizer interface that all auth methods must implement
type Authorizer interface {
    AddAuth(req *http.Request) error
}

// Built-in authorization methods
type Bearer struct {
    token string
}

type OAuth2ClientCredentials struct {
    config OAuth2Config
    tokenSource oauth2.TokenSource
}

// Usage examples:
// Bearer token
bearerAuth := auth.NewBearer("token")

// OAuth2 client credentials
oauth2Auth := auth.NewOAuth2ClientCredentials(auth.OAuth2Config{
    ClientID:     "client_id",
    ClientSecret: "client_secret",
    TokenURL:     "token_url",
    Scopes:       []string{"scope1", "scope2"},
})

// Custom authorization implementation
type CustomAuth struct {
    // Custom auth fields
}

func NewCustomAuth(...) *CustomAuth {
    // Initialize and return custom auth implementation
}

func (a *CustomAuth) AddAuth(req *http.Request) error {
    // custom auth logic

    req.Header.Set(header, value)

    return nil
}

// Initialize custom authentication with required parameters
customAuth := NewCustomAuth(...)
```

## Custom Events

```go
import (
    "github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

customEvent := secevent.EventType("https://example.com/events/custom")

// Implement Event interface for customEvent: https://github.com/SGNL-ai/caep.dev/tree/AddGoSsfSetLibrary/secevent#defining-custom-events

builder.New(transmitterURL,
    builder.WithEventTypes([]event.EventType{
        customEvent,
    }))

```

## Best Practices

1. **Authorization Management**
   - Use the default auth from the builder for most operations
   - Override auth only when necessary for specific operations

2. **Event Processing**
   - Use secevent's parser for proper SET validation and parsing

3. **Stream Management**
   - Always use defer for proper stream cleanup (Disable, Pause, or Delete)
   - Use the retry configuration for handling transient errors.

4. **Subject Management**
   - Use secevent's subject package for all subject operations

## Contributing

Contributions to the project are welcome, including feature enhancements, bug fixes, and documentation improvements.