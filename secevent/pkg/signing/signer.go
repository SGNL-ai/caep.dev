package signing

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Signer defines the interface for signing tokens
type Signer interface {
	Sign(ctx context.Context, claims jwt.Claims) (string, error)
}

// DefaultSigner uses the user supplied private key to sign
type DefaultSigner struct {
	signingKey    crypto.PrivateKey
	signingMethod jwt.SigningMethod
	keyID         *string
}

// SignerOption defines the function signature for signer options
type SignerOption func(*DefaultSigner)

func WithKeyID(kid string) SignerOption {
	return func(s *DefaultSigner) {
		s.keyID = &kid
	}
}

func WithSigningMethod(method jwt.SigningMethod) SignerOption {
	return func(s *DefaultSigner) {
		s.signingMethod = method
	}
}

func NewSigner(signingKey crypto.PrivateKey, opts ...SignerOption) (*DefaultSigner, error) {
	signer := &DefaultSigner{
		signingKey: signingKey,
	}

	switch k := signingKey.(type) {
	case *rsa.PrivateKey:
		signer.signingMethod = jwt.SigningMethodRS256
	case *ecdsa.PrivateKey:
		signer.signingMethod = jwt.SigningMethodES256
	default:
		return nil, fmt.Errorf("unsupported key type: %T", k)
	}

	for _, opt := range opts {
		opt(signer)
	}

	return signer, nil
}

func (s *DefaultSigner) Sign(ctx context.Context, claims jwt.Claims) (string, error) {
	// ctx is unused in this implementation since signing is purely local and operation is fast
	token := jwt.NewWithClaims(s.signingMethod, claims)

	if s.keyID != nil {
		token.Header["kid"] = *s.keyID
	}

	token.Header["typ"] = "secevent+jwt"

	return token.SignedString(s.signingKey)
}

func (s *DefaultSigner) SigningKey() crypto.PrivateKey {
	return s.signingKey
}

func (s *DefaultSigner) SigningMethod() jwt.SigningMethod {
	return s.signingMethod
}

func (s *DefaultSigner) KeyID() *string {
	return s.keyID
}
