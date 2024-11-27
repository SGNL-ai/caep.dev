package parser

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/sgnl-ai/caep.dev/secevent/pkg/builder"

	_ "github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/caep" // Initialize CAEP events
	_ "github.com/sgnl-ai/caep.dev/secevent/pkg/schemes/ssf"  // Initialize SSF events
)

// Parser parses and validates SecEvents
type Parser struct {
	keySet           jwk.Set
	jwksURL          *url.URL
	expectedIssuer   string
	expectedAudience []string
}

// Option defines the function signature for parser options
type Option func(*Parser)

// WithJWKSURL sets a JWKS URL to fetch keys from
func WithJWKSURL(rawURL string) Option {
	return func(p *Parser) {
		if parsedURL, err := url.Parse(rawURL); err == nil {
			p.jwksURL = parsedURL
			p.keySet = nil // Clear any existing key set
		}
	}
}

// WithJWKSURLParsed sets a JWKS URL using a pre-parsed URL
func WithJWKSURLParsed(parsedURL *url.URL) Option {
	return func(p *Parser) {
		p.jwksURL = parsedURL
		p.keySet = nil // Clear any existing key set
	}
}

// WithJWKSJSON sets the JWKS JSON data
func WithJWKSJSON(jwksJSON []byte) Option {
	return func(p *Parser) {
		if keySet, err := jwk.Parse(jwksJSON); err == nil {
			p.keySet = keySet
			p.jwksURL = nil // Clear any existing JWKS URL
		}
	}
}

// WithPublicKey sets a direct public key
func WithPublicKey(key interface{}, kid string) Option {
	return func(p *Parser) {
		if rawKey, err := jwk.FromRaw(key); err == nil {
			// Set the key ID
			if err := rawKey.Set(jwk.KeyIDKey, kid); err != nil {
				return
			}

			keySet := jwk.NewSet()
			if err := keySet.AddKey(rawKey); err != nil {
				return
			}

			p.keySet = keySet
			p.jwksURL = nil // Clear any existing JWKS URL
		}
	}
}

// WithExpectedIssuer sets the expected issuer for validation
func WithExpectedIssuer(issuer string) Option {
	return func(p *Parser) {
		p.expectedIssuer = issuer
	}
}

// WithExpectedAudience sets the expected audience for validation
func WithExpectedAudience(audience ...string) Option {
	return func(p *Parser) {
		p.expectedAudience = audience
	}
}

// NewParser creates a new SecEvent parser with the provided options
func NewParser(opts ...Option) *Parser {
	p := &Parser{}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

// fetchJWKS fetches the JWKS from the configured URL
func (p *Parser) fetchJWKS() (jwk.Set, error) {
	if p.jwksURL == nil {
		return nil, fmt.Errorf("no JWKS URL configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	keySet, err := jwk.Fetch(ctx, p.jwksURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	return keySet, nil
}

// getKey returns the appropriate key for JWT verification
func (p *Parser) getKey(token *jwt.Token) (interface{}, error) {
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("token header does not contain 'kid'")
	}

	var keySet jwk.Set
	var err error

	// If we have a JWKS URL, fetch fresh keys
	if p.jwksURL != nil {
		keySet, err = p.fetchJWKS()
		if err != nil {
			return nil, err
		}
	} else if p.keySet != nil {
		keySet = p.keySet
	} else {
		return nil, fmt.Errorf("no keys available for verification")
	}

	key, found := keySet.LookupKeyID(kid)
	if !found {
		return nil, fmt.Errorf("no key found for kid %s", kid)
	}

	var rawKey interface{}
	if err := key.Raw(&rawKey); err != nil {
		return nil, fmt.Errorf("failed to get raw key: %w", err)
	}

	return rawKey, nil
}

// getParserOptions returns the appropriate JWT parser options based on configuration
func (p *Parser) getParserOptions() []jwt.ParserOption {
	options := []jwt.ParserOption{
		jwt.WithValidMethods([]string{"RS256", "ES256"}),
	}

	if p.expectedIssuer != "" {
		options = append(options, jwt.WithIssuer(p.expectedIssuer))
	}

	// Handle audience validation
	if len(p.expectedAudience) > 0 {
		for _, aud := range p.expectedAudience {
			options = append(options, jwt.WithAudience(aud))
		}
	}

	return options
}

// ParseMultiSecEvent parses and validates a signed SecEvent
func (p *Parser) ParseMultiSecEvent(tokenString string) (*builder.MultiSecEvent, error) {
	var set builder.MultiSecEvent

	parser := jwt.NewParser(p.getParserOptions()...)

	token, err := parser.ParseWithClaims(tokenString, &set, p.getKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &set, nil
}

// ParseSecEvent parses and validates a signed SingleEventSecEvent
func (p *Parser) ParseSecEvent(tokenString string) (*builder.SecEvent, error) {
	var set builder.SecEvent

	parser := jwt.NewParser(p.getParserOptions()...)

	token, err := parser.ParseWithClaims(tokenString, &set, p.getKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &set, nil
}

// ParseMultiSecEventNoVerify parses a SecEvent without validation
func (p *Parser) ParseMultiSecEventNoVerify(tokenString string) (*builder.MultiSecEvent, error) {
	var set builder.MultiSecEvent

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	_, _, err := parser.ParseUnverified(tokenString, &set)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return &set, nil
}

// ParseSecEventNoVerify parses a SingleEventSecEvent without validation
func (p *Parser) ParseSecEventNoVerify(tokenString string) (*builder.SecEvent, error) {
	var set builder.SecEvent

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	_, _, err := parser.ParseUnverified(tokenString, &set)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return &set, nil
}
