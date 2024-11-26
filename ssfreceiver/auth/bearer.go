package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// BearerAuth implements bearer token authorization
type BearerAuth struct {
	token      string
	tokenMutex sync.RWMutex
}

func NewBearer(token string) (*BearerAuth, error) {
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	auth := &BearerAuth{
		token: token,
	}

	return auth, nil
}

// AddAuth implements the Authorizer interface
func (b *BearerAuth) AddAuth(ctx context.Context, req *http.Request) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	token := b.getToken()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return nil
}

func (b *BearerAuth) SetToken(token string) error {
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	b.tokenMutex.Lock()
	defer b.tokenMutex.Unlock()

	b.token = token

	return nil
}

func (b *BearerAuth) getToken() string {
	b.tokenMutex.RLock()
	defer b.tokenMutex.RUnlock()

	return b.token
}
