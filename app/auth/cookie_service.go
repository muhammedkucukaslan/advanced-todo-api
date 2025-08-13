package auth

import (
	"context"
	"time"
)

type CookieService interface {
	SetRefreshToken(ctx context.Context, claims *RefreshTokenCookieClaims)
}

type RefreshTokenCookieClaims struct {
	Token    string
	Duration time.Duration
	Secure   bool
}
