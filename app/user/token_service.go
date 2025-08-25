package user

type TokenService interface {
	GenerateSecureEmailToken(email string) (string, error)
	ValidateSecureEmailToken(tokenString string) (string, error)
}
