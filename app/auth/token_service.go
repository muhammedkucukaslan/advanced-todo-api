package auth

type TokenPayload struct {
	UserID string `json:"userID"`
	Role   string `json:"role"`
}

type TokenService interface {
	GenerateAuthToken(userID string, role string) (string, error)
	ValidateAuthToken(token string) (TokenPayload, error)
	GenerateVerificationToken(email string) (string, error)
}
