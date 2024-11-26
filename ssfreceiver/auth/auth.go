package auth

import (
	"context"
	"net/http"
)

// Authorizer is the interface that authorization method must implement
type Authorizer interface {
	// AddAuth adds authorization to the provided request
	AddAuth(ctx context.Context, req *http.Request) error
}
