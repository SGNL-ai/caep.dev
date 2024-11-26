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
  - [Push Event Reception](#push-event-reception)
- [Subject Management](#subject-management)
- [Authorization](#authorization)
- [Types and Constants](#types-and-constants)
  - [Stream Status Types](#stream-status-types)
  - [Delivery Methods](#delivery-methods)
- [Builder Options](#builder-options)
- [Operation Options](#operation-options)
- [Retry Configuration](#retry-configuration)
- [Stream Interface](#stream-interface)
- [Custom Events](#custom-events)
- [Best Practices](#best-practices)
- [Contributing](#contributing)

## Features

- **Complete SSF Receiver Implementation**: Full support for creating and managing SSF streams according to the OpenID SSF specification
- **Integrates with secevent**: Integrates with the [secevent](https://github.com/SGNL-ai/caep.dev/tree/AddGoSsfSetLibrary/secevent) library for SET handling, subject management, and event parsing
- **Multiple Delivery Methods**: Support for both poll-based and push-based event delivery
- **Flexible Authorization**: Support for various built-in authorization methods and custom implementation.
- **Stream Lifecycle Management**: Comprehensive stream status control (enable, disable, pause)
- **Configurable Retry Mechanism**: Support for configurable retry mechanism for handling transient failures
- **Customizable Headers**: Support for endpoint-specific headers

## Dependencies

This library requires:
- [secevent](https://github.com/SGNL-ai/caep.dev/tree/AddGoSsfSetLibrary/secevent)

## Installation

```bash
go get github.com/caep.dev-receiver/ssfreceiver
```

## Quick Start

### Authorization Setup

First, set up your authorization method:

```go
package main

import (
    "github.com/caep.dev-receiver/ssfreceiver/auth"
    "golang.org/x/oauth2/clientcredentials"
)

// Bearer token authentication
bearerAuth, err := auth.NewBearer("token")
if err != nil {
    // Handle error
}

// OAuth2 client credentials
oauth2Auth, err := auth.NewOAuth2ClientCredentials(&clientcredentials.Config{
    ClientID:     "client_id",
    ClientSecret: "client_secret",
    TokenURL:     "token_url",
    Scopes:       []string{"scope1", "scope2"},
})
if err != nil {
    // Handle error
}

// Custom authorization implementation
type CustomAuth struct {
    // Custom auth fields
}

// AddAuth implements the auth.Authorizer interface
func (a *CustomAuth) AddAuth(ctx context.Context, req *http.Request) error {
    // Custom auth logic
    req.Header.Set("Authorization", "Custom scheme-name")
    
    return nil
}

// Initialize custom authentication with required parameters
customAuth := &CustomAuth{}
```

### Poll-Based Stream

```go
package main

import (
    "context"
    "log"
    "github.com/caep.dev-receiver/ssfreceiver/builder"
    "github.com/caep.dev-receiver/ssfreceiver/options"
    "github.com/sgnl-ai/caep.dev-receiver/secevent/pkg/subject"
    "github.com/sgnl-ai/caep.dev-receiver/secevent/pkg/schemes/caep"
    "github.com/sgnl-ai/caep.dev-receiver/secevent/pkg/event"
)

func main() {
    // Create a stream builder
    streamBuilder, err := builder.New(
        "https://transmitter.example.com/.well-known/ssf-configuration",
        builder.WithPollDelivery(),
        builder.WithAuth(bearerAuth),
        builder.WithEventTypes([]event.EventType{
            caep.EventTypeSessionRevoked,      
            caep.EventTypeTokenClaimsChange,
            event.EventType("https://custom.example.com/events/custom"),
        }),
        builder.WithExistingCheck(),              // Optional, checks for existing streams
    )
    if err != nil {
        // Handle error
    }

    // Setup stream
    ssfStream, err := streamBuilder.Setup(context.Background())
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
        options.WithAuth(customAuth), // Optional, overrides stream's default auth for this request
    )
    if err != nil {
        // Handle error
    }

    // Poll for events with all available options
    events, err := ssfStream.Poll(context.Background(),
        options.WithMaxEvents(10),          // Optional, default is 100
        options.WithAutoAck(true),          // Optional, immediately auto acknowledges received events
        options.WithAckJTIs([]string{"jti-1", "jti-2"}), // Optional, acknowledges specified events
        options.WithAuth(customAuth),       // Optional, overrides stream's default auth for this request
        options.WithHeaders(map[string]string{  // Optional, adds additional headers or overrides headers for this request
            "Custom-Header": "value",
        }),
    )
    if err != nil {
        // Handle error
    }

    // Process events using secevent parser
    secEventParser := parser.NewParser(
        parser.WithJWKSURL("https://issuer.example.com/jwks.json"),
        parser.WithExpectedIssuer("https://issuer.example.com"),
        parser.WithExpectedAudience("https://receiver.example.com"),
    )

    for _, rawEvent := range events {
        parsedEvent, err := secEventParser.ParseSecEvent(rawEvent)
        if err != nil {
            // Handle error
            continue
        }
        
        handleEvent(parsedEvent)
    }
}
```

### Push-Based Stream

```go
package main

import (
    "context"
    "net/http"
    "github.com/caep.dev-receiver/ssfreceiver/builder"
    "github.com/sgnl-ai/caep.dev-receiver/secevent/pkg/parser"
    "github.com/sgnl-ai/caep.dev-receiver/secevent/pkg/schemes/caep"
)

func main() {
    // Create a stream builder
    streamBuilder, err := builder.New(
        "https://transmitter.example.com/.well-known/ssf-configuration",
        builder.WithPushDelivery("https://receiver.example.com/events"),
        builder.WithAuth(bearerAuth),
        builder.WithEventTypes([]event.EventType{
            caep.EventTypeSessionRevoked,      
            caep.EventTypeTokenClaimsChange,
            event.EventType("https://custom.example.com/events/custom"),
        })
    )
    if err != nil {
        // Handle error
    }

    // Build and validate the stream
    ssfStream, err := streamBuilder.Setup(context.Background())
    if err != nil {
        // Handle error
    }

    defer ssfStream.Disable(context.Background())
}
```

```go
// Use secevent parser to parse incoming events
secEventParser := parser.NewParser(
    parser.WithJWKSURL("https://issuer.example.com/jwks.json"),
    parser.WithExpectedIssuer("https://issuer.example.com"),
    parser.WithExpectedAudience("https://receiver.example.com"),
)

// Push-endpoint(https://receiver.example.com/events) handler
http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
    event, err := secEventParser.ParseSecEvent(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        
        return
    }
    
    handleEvent(event)

    w.WriteHeader(http.StatusOK)
})
```

## Stream Management

### Stream Creation
```go
// Create builder with all available options
builder, err := builder.New(transmitterURL,
    // Delivery method options
    builder.WithPollDelivery(),  // or WithPushDelivery(endpoint)
    
    // Authentication
    builder.WithAuth(auth.NewBearer("default-token")),
    
    // Event configuration
    builder.WithEventTypes([]event.EventType{
        caep.EventTypeSessionRevoked,
        caep.EventTypeTokenClaimsChange,
    }),
    
    // Stream handling
    builder.WithExistingCheck(),  // Optional, checks for existing streams
    builder.WithDescription("Stream description"),  // Optional, adds stream description
    
    // HTTP client configuration
    builder.WithHTTPClient(&http.Client{    // Optional, uses user-supplied http client
        Timeout: time.Second * 30,
    }),
    
    // Retry configuration
    builder.WithRetryConfig(config.RetryConfig{ // Optional, if not supplied, uses default retry config
        MaxRetries:        3,
        InitialBackoff:    time.Second,
        MaxBackoff:        time.Second * 30,
        BackoffMultiplier: 2.0,
        RetryableStatus:   map[int]bool{
            408: true, // Request Timeout
            429: true, // Too Many Requests
            500: true, // Internal Server Error
            502: true, // Bad Gateway
            503: true, // Service Unavailable
            504: true, // Gateway Timeout
        },
    }),
    
    // Endpoint-specific headers
    builder.WithMetadataEndpointHeaders(map[string]string{  // Optional, adds additional headers or overrides headers for this endpoint
        "Custom-Metadata-Header": "value",
    }),
    builder.WithConfigurationEndpointHeaders(map[string]string{ // Optional, adds additional headers or overrides headers for this endpoint
        "Custom-Config-Header": "value",
    }),
    builder.WithStatusEndpointHeaders(map[string]string{    // Optional, adds additional headers or overrides headers for this endpoint
        "Custom-Status-Header": "value",
    }),
    builder.WithAddSubjectEndpointHeaders(map[string]string{    // Optional, adds additional headers or overrides headers for this endpoint
        "Custom-AddSubject-Header": "value",
    }),
    builder.WithRemoveSubjectEndpointHeaders(map[string]string{ // Optional, adds additional headers or overrides headers for this endpoint
        "Custom-RemoveSubject-Header": "value",
    }),
    builder.WithVerificationEndpointHeaders(map[string]string{  // Optional, adds additional headers or overrides headers for this endpoint
        "Custom-Verify-Header": "value",
    }),
    
    // Global headers for all endpoints
    builder.WithEndpointHeaders(map[string]string{  // Optional, adds additional headers or overrides headers for all endpoints
        "Global-Custom-Header": "value",
    }),
)
if err != nil {
    // Handle error
}

// Build stream
stream, err := builder.Setup(ctx)
if err != nil {
    // Handle error
}
```

### Stream Configuration
```go
// Get current configuration
config, err := stream.GetConfiguration(ctx,
    options.WithAuth(customAuth),  // Optional, overrides stream's default auth for this request
    options.WithHeaders(map[string]string{ // Optional, adds additional headers or overrides headers for this request
        "Custom-Header": "value",
    }),
)
if err != nil {
    // Handle error
}

// Update stream configuration
updatedConfig := &types.StreamConfigurationRequest{
    StreamID: stream.GetStreamID(),
    Delivery: &types.DeliveryConfig{
        Method:      types.DeliveryMethodPush,
        EndpointURL: endpointURL,
    },
    EventsRequested: []event.EventType{
        caep.EventTypeSessionRevoked,
        caep.EventTypeTokenClaimsChange,
    },
    Description: "Updated configuration",
}

newConfig, err := stream.UpdateConfiguration(ctx, 
    updatedConfig,
    options.WithAuth(customAuth), // Optional, overrides stream's default auth for this request
    options.WithHeaders(map[string]string{ // Optional, adds additional headers or overrides headers for this request
        "Custom-Header": "value",
    }),
)
```

### Stream Status
```go
// Get current status
status, err := stream.GetStatus(ctx,
    options.WithAuth(customAuth), // Optional, overrides stream's default auth for this request
    options.WithHeaders(map[string]string{ // Optional, adds additional headers or overrides headers for this request
        "Custom-Header": "value",
    }),
)
if err != nil {
    // Handle error
}

// Pause stream
err = stream.Pause(ctx, 
    options.WithStatusReason("System maintenance"), // Optional, specifies the reason for pausing the stream
    options.WithAuth(customAuth), // Optional, overrides stream's default auth for this request
    options.WithHeaders(map[string]string{  // Optional, adds additional headers or overrides headers for this request
        "Custom-Header": "value",
    }),
)

// Disable stream
err = stream.Disable(ctx,
    options.WithStatusReason("Disabling stream"),   // Optional, specifies the reason for disabling the stream
    options.WithAuth(customAuth),   // Optional, overrides stream's default auth for this request
)

// Delete stream
err = stream.Delete(ctx,
    options.WithAuth(customAuth),   // Optional, overrides stream's default auth for this request
    options.WithHeaders(map[string]string{   // Optional, adds additional headers or overrides headers for this request
        "Custom-Header": "value",
    }),
)

// Resume stream
err = stream.Resume(ctx,
    options.WithStatusReason("Resuming normal operation"),  // Optional, specifies the reason for resuming the stream
    options.WithAuth(customAuth),   // Optional, overrides stream's default auth for this request
)

// Update status directly
err = stream.UpdateStatus(ctx,
    types.StatusEnabled,
    options.WithStatusReason("Custom status update"),   // Optional, specifies the reason for updating the stream
    options.WithAuth(customAuth),   // Optional, overrides stream's default auth for this request
)
```

### Stream Verification
```go
// Verify stream with all available options
err := stream.Verify(ctx,
    options.WithState("verification-state"),    // Optional, adds the verification event's state
    options.WithAuth(customAuth),   // Optional, overrides stream's default auth for this request
    options.WithHeaders(map[string]string{  // Optional, adds additional headers or overrides headers for this request
        "Custom-Header": "value",
    }),
)
if err != nil {
    // Handle error
}
```

## Event Handling

### Polling Events
```go
// Poll with all available options
events, err := stream.Poll(ctx,
    options.WithMaxEvents(10),          // Optional, default is 100
    options.WithAutoAck(true),          // Optional, immediately auto acknowledges the polled events
    options.WithAckJTIs([]string{...}), // Optional, acknowledges previouly polled events
    options.WithAuth(customAuth),       // Optional, overrides stream's default auth for this request
    options.WithHeaders(map[string]string{  // Optional, adds additional headers or overrides headers for this request
        "Custom-Header": "value",
    }),
)
if err != nil {
    // Handle error
}

// Process events using secevent parser
secEventParser := secevent.NewParser(
    parser.WithJWKSURL("https://issuer.example.com/jwks.json"),
    parser.WithExpectedIssuer("https://issuer.example.com"),
)

for _, rawEvent := range events {
    parsedEvent, err := secEventParser.ParseSecEvent(rawEvent)
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
// Acknowledge specific events with all available options
err := stream.Acknowledge(ctx, 
    []string{"jti1", "jti2"},
    options.WithAuth(customAuth),   // Optional, overrides stream's default auth for this request
    options.WithHeaders(map[string]string{  // Optional, adds additional headers or overrides headers for this request
        "Custom-Header": "value",
    }),
)
if err != nil {
    // Handle error
}
```

### Push Event Reception
```go
// Set up secevent parser for handling incoming events
secEventParser := secevent.NewParser(
    parser.WithJWKSURL("https://issuer.example.com/jwks.json"),
    parser.WithExpectedIssuer("https://issuer.example.com"),
    parser.WithExpectedAudience("https://receiver.example.com"),
)

// Set up HTTP handler for receiving push events
http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
    // Parse and validate the incoming SET
    event, err := secEventParser.ParseSecEvent(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Handle the event based on type
    switch event.Type() {
    case caep.EventTypeSessionRevoked:
        if err := handleSessionRevoked(event); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    case caep.EventTypeTokenClaimsChange:
        if err := handleTokenClaimsChange(event); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    default:
        if err := handleUnknownEvent(event); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    w.WriteHeader(http.StatusOK)
})
```

## Subject Management

Using secevent's subject package for subject creation and management:

```go
import "github.com/sgnl-ai/caep.dev-receiver/secevent/pkg/subject"

// Create various subject types using secevent
emailSubject, err := subject.NewEmailSubject("user@example.com")
if err != nil {
    // Handle error
}

phoneSubject, err := subject.NewPhoneSubject("+1-555-123-4567")
if err != nil {
    // Handle error
}

issSubSubject, err := subject.NewIssSubSubject("https://issuer.example.com", "user123")
if err != nil {
    // Handle error
}

// Complex subject with multiple identifiers
complexSubject := subject.NewComplexSubject().
    WithUser(emailSubject).
    WithDevice(subject.NewOpaqueSubject("device-123"))

// Add subject to stream with all available options
err = stream.AddSubject(ctx,
    emailSubject,
    options.WithSubjectVerification(true),  // Optional, sets subject verification status
    options.WithAuth(customAuth),   // Optional, overrides stream's default auth for this request
    options.WithHeaders(map[string]string{  // Optional, adds additional headers or overrides headers for this request
        "Custom-Header": "value",
    }),
)
if err != nil {
    // Handle error
}

// Remove subject from stream with all available options
err = stream.RemoveSubject(ctx,
    emailSubject,
    options.WithAuth(customAuth),   // Optional, overrides stream's default auth for this request
    options.WithHeaders(map[string]string{  // Optional, adds additional headers or overrides headers for this request
        "Custom-Header": "value",
    }),
)
if err != nil {
    // Handle error
}
```

## Authorization

The library provides a flexible authorization system through the `auth` package. Both built-in authorization methods and custom implementations are supported.

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "golang.org/x/oauth2/clientcredentials"
    "github.com/caep.dev-receiver/ssfreceiver/auth"
)

// 1. Built-in Bearer Token Authentication
bearerAuth, err := auth.NewBearer("your-token")
if err != nil {
    // Handle error
}

// Update bearer token if needed
err = bearerAuth.SetToken("new-token")
if err != nil {
    // Handle error
}

// 2. Built-in OAuth2 Client Credentials Authentication
oauth2Auth, err := auth.NewOAuth2ClientCredentials(&clientcredentials.Config{
    ClientID:     "client_id",
    ClientSecret: "client_secret",
    TokenURL:     "https://auth.example.com/token",
    Scopes:       []string{"scope1", "scope2"},
})
if err != nil {
    // Handle error
}

// 3. Custom Authorization Implementation
// First, create your custom authorizer type
type CustomAuth struct {
    apiKey string
    // Add any other required fields
}

// Implement the Authorizer interface
func (a *CustomAuth) AddAuth(ctx context.Context, req *http.Request) error {
    if req == nil {
        return fmt.Errorf("request cannot be nil")
    }
    
    // Add your custom authorization logic
    req.Header.Set("X-API-Key", a.apiKey)
    // Add any other headers or auth-related modifications
    
    return nil
}

// Create an instance of your custom authorizer
customAuth := &CustomAuth{
    apiKey: "your-api-key",
}

// Use any of these authorizers with the stream builder
streamBuilder, err := builder.New(
    "https://transmitter.example.com/.well-known/ssf-configuration",
    builder.WithAuth(bearerAuth), // or oauth2Auth or customAuth
)
```

## Types and Constants

### Stream Status Types
```go
// Possible stream status values
const (
    StatusEnabled  StreamStatusType = "enabled"   // Transmitter must transmit events
    StatusPaused   StreamStatusType = "paused"    // Transmitter must not transmit but will hold events
    StatusDisabled StreamStatusType = "disabled"  // Transmitter must not transmit and will not hold events
)
```

### Delivery Methods
```go
// Available delivery methods
const (
    DeliveryMethodPush DeliveryMethod = "urn:ietf:rfc:8935"  // Push-based delivery
    DeliveryMethodPoll DeliveryMethod = "urn:ietf:rfc:8936"  // Poll-based delivery
)
```

## Builder Options

All available options for stream builder configuration:

```go
// Delivery Method Options. One of them should be supplied
WithPollDelivery()                    // Configure for poll-based delivery
WithPushDelivery(pushEndpoint string) // Configure for push-based delivery

// Authentication
WithAuth(auth.Authorizer)  // Set default authorization method

// Event Configuration
WithEventTypes([]event.EventType)     // Set event types to receive

// Stream Management
WithExistingCheck()                   // Optional. Enable checking for existing streams
WithDescription(string)               // Optional. Set stream description

// HTTP Configuration
WithHTTPClient(*http.Client)          // Optional. Set custom HTTP client

// Retry Configuration
WithRetryConfig(config.RetryConfig)   // Optional. Configure retry behavior

// Endpoint-Specific Headers
WithMetadataEndpointHeaders(map[string]string)       // Optional. Additional headers for metadata endpoint
WithConfigurationEndpointHeaders(map[string]string)  // Optional. Additional headers for configuration endpoint
WithStatusEndpointHeaders(map[string]string)         // Optional. Additional headers for status endpoint
WithAddSubjectEndpointHeaders(map[string]string)     // Optional. Additional headers for add subject endpoint
WithRemoveSubjectEndpointHeaders(map[string]string)  // Optional. Additional headers for remove subject endpoint
WithVerificationEndpointHeaders(map[string]string)   // Optional. Additional headers for verification endpoint

// Global Headers
WithEndpointHeaders(map[string]string)               // Optional. Additional headers for all endpoints
```

## Operation Options

All available options for stream operations:

```go
// Authentication Options
WithAuth(auth.Authorizer)           // Optional. Override default authorizer

// Event Handling Options
WithMaxEvents(int)                  // Optional. Set maximum events to poll (default: 100)
WithAutoAck(bool)                   // Optional. Enable automatic acknowledgment of polled events
WithAckJTIs([]string)              // Optional. Specify JTIs to acknowledge in a poll operation

// Verification Options
WithState(string)                   // Optional. Set state for verification
WithSubjectVerification(bool)       // Optional. Set subject verification status

// HTTP Options
WithHeaders(map[string]string)      // Optional. Set additional headers for the operation

// Status Options
WithStatusReason(string)            // Optional. Set reason for status change
```

## Retry Configuration

Configuration options for retry behavior:

```go
type RetryConfig struct {
    // MaxRetries is the maximum number of retry attempts (default: 3)
    MaxRetries int
    
    // InitialBackoff is the initial delay between retry attempts (default: 1 second)
    InitialBackoff time.Duration
    
    // MaxBackoff is the maximum delay between retry attempts (default: 30 seconds)
    MaxBackoff time.Duration
    
    // BackoffMultiplier is the factor by which the backoff increases (default: 2.0)
    BackoffMultiplier float64
    
    // RetryableStatus is a map of HTTP status codes that should trigger a retry
    RetryableStatus map[int]bool  // Default: 408, 429, 500, 502, 503, 504
}

// Default retry configuration
config.DefaultRetryConfig() // Returns default configuration
```

## Stream Interface

Complete interface for stream operations:

```go
type Stream interface {
    // Get stream ID
    GetStreamID() string

    // Get transmitter metadata
    GetMetadata() *types.TransmitterMetadata
    
    // Configuration management
    GetConfiguration(ctx context.Context, opts ...options.Option) (*types.StreamConfiguration, error)
    UpdateConfiguration(ctx context.Context, config *types.StreamConfigurationRequest, opts ...options.Option) (*types.StreamConfiguration, error)
    
    // Status management
    GetStatus(ctx context.Context, opts ...options.Option) (*types.StreamStatus, error)
    UpdateStatus(ctx context.Context, status types.StreamStatusType, opts ...options.Option) error
    
    // Subject management
    AddSubject(ctx context.Context, sub subject.Subject, opts ...options.Option) error
    RemoveSubject(ctx context.Context, sub subject.Subject, opts ...options.Option) error
    
    // Stream verification
    Verify(ctx context.Context, opts ...options.Option) error
    
    // Event handling
    Poll(ctx context.Context, opts ...options.Option) (map[string]string, error)
    Acknowledge(ctx context.Context, jtis []string, opts ...options.Option) error
    
    // Stream lifecycle
    Delete(ctx context.Context, opts ...options.Option) error
    Pause(ctx context.Context, opts ...options.Option) error
    Resume(ctx context.Context, opts ...options.Option) error
    Disable(ctx context.Context, opts ...options.Option) error
}
```

## Custom Events

```go
import (
    "github.com/sgnl-ai/caep.dev/secevent/pkg/event"
)

customEvent := secevent.EventType("https://example.com/events/custom")

// Implement Event interface for customEvent: https://github.com/SGNL-ai/caep.dev/tree/main/secevent#defining-custom-events

builder.New(transmitterURL,
    builder.WithEventTypes([]event.EventType{
        customEvent, // Request the custom event
    }))

```

Note: Custom events must follow the SET event specification and should use unique URIs to avoid conflicts with standard event types.

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