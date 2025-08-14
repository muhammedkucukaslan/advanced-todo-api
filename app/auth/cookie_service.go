package auth

import (
	"context"
	"time"
)

type CookieService interface {
	SetRefreshToken(ctx context.Context, claims *RefreshTokenCookieClaims)
	RemoveRefreshToken(ctx context.Context)
}

type RefreshTokenCookieClaims struct {
	Token    string
	Duration time.Duration
	Secure   bool
}

// TODO add logout handler and delete refresh token from db when user logs out
