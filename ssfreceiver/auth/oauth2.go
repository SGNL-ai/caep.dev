package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// OAuth2Auth implements OAuth2 client credentials authorization
type OAuth2Auth struct {
	config       *clientcredentials.Config
	currentToken *oauth2.Token
	tokenMutex   sync.RWMutex
}

func NewOAuth2ClientCredentials(config *clientcredentials.Config) (*OAuth2Auth, error) {
	if config.TokenURL == "" {
		return nil, fmt.Errorf("token URL is required")
	}

	if config.ClientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}

	if config.ClientSecret == "" {
		return nil, fmt.Errorf("client secret is required")
	}

	return &OAuth2Auth{
		config: config,
	}, nil
}

// AddAuth implements the Authorizer interface
func (a *OAuth2Auth) AddAuth(ctx context.Context, req *http.Request) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	token, err := a.getValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get valid token: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	return nil
}

func (a *OAuth2Auth) getValidToken(ctx context.Context) (*oauth2.Token, error) {
	a.tokenMutex.Lock()
	defer a.tokenMutex.Unlock()

	if a.currentToken == nil {
		return a.fetchNewToken(ctx)
	}

	if !a.currentToken.Valid() {
		return a.fetchNewToken(ctx)
	}

	return a.currentToken, nil
}

func (a *OAuth2Auth) fetchNewToken(ctx context.Context) (*oauth2.Token, error) {
	token, err := a.config.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	a.currentToken = token

	return token, nil
}
