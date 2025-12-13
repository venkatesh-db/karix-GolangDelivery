package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT payload propagated through the system.
type Claims struct {
	TenantID string   `json:"tenant_id"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// Validator verifies bearer tokens using a shared secret.
type Validator struct {
	secret []byte
}

// NewValidator returns a Validator using the provided HMAC secret.
func NewValidator(secret string) *Validator {
	return &Validator{secret: []byte(secret)}
}

var (
	ErrMissingToken = errors.New("authorization token missing")
	ErrInvalidToken = errors.New("authorization token invalid")
)

// Parse reads the Authorization header value (Bearer token) and returns claims.
func (v *Validator) Parse(header string) (*Claims, error) {
	tokenString := extractBearer(header)
	if tokenString == "" {
		return nil, ErrMissingToken
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return v.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.TenantID == "" {
		return nil, ErrInvalidToken
	}
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func extractBearer(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// Context helpers

type claimsKey struct{}

// WithClaims adds claims to context.
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

// ClaimsFromContext retrieves claims.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*Claims)
	return claims, ok && claims != nil
}
