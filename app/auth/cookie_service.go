package auth

import (
	"context"
)

type CookieService interface {
	SetRefreshToken(ctx context.Context, token string)
	SetAccessToken(ctx context.Context, token string)
	RemoveTokens(ctx context.Context)
}
