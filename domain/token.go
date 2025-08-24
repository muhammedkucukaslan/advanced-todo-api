package domain

import (
	"time"

	"github.com/google/uuid"
)

var (
	RefreshTokenCookieName   = "refresh_token"
	RefreshCookiePath        = "/auth"
	RefreshTokenCookieMaxAge = 30 * 24 * time.Hour
	AccessTokenCookieMaxAge  = 15 * time.Minute
	AccessTokenCookiePath    = "/"
	AccessTokenCookieName    = "access_token"
	CookieSecure             = true
)

type RefreshToken struct {
	Id        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func NewRefreshToken(userID uuid.UUID, token string) *RefreshToken {
	return &RefreshToken{
		Id:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(RefreshTokenCookieMaxAge),
		CreatedAt: time.Now(),
	}
}
