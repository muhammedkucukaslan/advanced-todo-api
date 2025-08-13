package user

type TokenService interface {
	GenerateTokenForForgotPassword(email string) (string, error)
	ValidateForgotPasswordToken(tokenString string) (string, error)
	ValidateVerifyEmailToken(tokenString string) (string, error)
	GenerateEmailVerificationToken(email string) (string, error)
}
