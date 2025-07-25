package auth

import "time"

type TokenPayload struct {
	UserID string `json:"userID"`
	Role   string `json:"role"`
}

type TokenService interface {
	GenerateToken(userID string, role string, createdAt time.Time) (string, error)
	ValidateToken(token string) (TokenPayload, error)
	GenerateVerificationToken(email string) (string, error)
}
