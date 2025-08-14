package auth

type TokenPayload struct {
	UserID string `json:"userID"`
	Role   string `json:"role"`
}

type TokenService interface {
	GenerateAuthAccessToken(userID string, role string) (string, error)
	ValidateAuthAccessToken(token string) (*TokenPayload, error)
	GenerateEmailVerificationToken(email string) (string, error)
	GenerateAuthRefreshToken(userID string, role string) (string, error)
}
