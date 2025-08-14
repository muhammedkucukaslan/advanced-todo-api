package domain

import (
	"time"

	"github.com/google/uuid"
)

var (
	RefreshTokenCookieName = "refresh_token"
)

type RefreshToken struct {
	Id        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func NewRefreshToken(userID uuid.UUID, token string, duration time.Duration) *RefreshToken {
	return &RefreshToken{
		Id:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(duration),
		CreatedAt: time.Now(),
	}
}
