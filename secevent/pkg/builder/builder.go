package builder

import (
	"github.com/sgnl-ai/caep.dev/secevent/pkg/id"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/token"
)

// Builder provides configuration for creating SecEvents
type Builder struct {
	defaultIssuer      string
	defaultIDGenerator id.Generator
}

// Option defines the function signature for builder options
type Option func(*Builder)

// WithDefaultIssuer sets the default issuer for all SecEvents created by this builder
func WithDefaultIssuer(issuer string) Option {
	return func(b *Builder) {
		b.defaultIssuer = issuer
	}
}

// WithDefaultIDGenerator sets the default ID generator for all SecEvents created by this builder
func WithDefaultIDGenerator(generator id.Generator) Option {
	return func(b *Builder) {
		b.defaultIDGenerator = generator
	}
}

// NewBuilder creates a new SecEvent builder with the provided options
func NewBuilder(opts ...Option) *Builder {
	b := &Builder{
		defaultIDGenerator: id.NewUUIDGenerator(), // Default to UUID generator
	}

	// Apply options
	for _, opt := range opts {
		opt(b)
	}

	return b
}

// NewMultiSecEvent creates a new multi-event SecEvent with default configurations
func (b *Builder) NewMultiSecEvent() *token.MultiSecEvent {
	secEvent := token.NewMultiSecEvent()
	if b.defaultIssuer != "" {
		secEvent.WithIssuer(b.defaultIssuer)
	}

	if b.defaultIDGenerator != nil {
		secEvent.WithID(b.defaultIDGenerator.Generate())
	}

	return secEvent
}

// NewSecEvent creates a new single-event SecEvent with default configurations
func (b *Builder) NewSecEvent() *token.SecEvent {
	secEvent := token.NewSecEvent()
	if b.defaultIssuer != "" {
		secEvent.WithIssuer(b.defaultIssuer)
	}

	if b.defaultIDGenerator != nil {
		secEvent.WithID(b.defaultIDGenerator.Generate())
	}

	return secEvent
}
