package jwt

import (
	"fmt"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"

	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	secretKey              string
	tokenDuration          time.Duration
	emailVerificationToken time.Duration
	forgotPasswordToken    time.Duration
}

func NewTokenService(secretKey string, tokenDuration, emailVerificationToken, forgotPasswordToken time.Duration) *TokenService {
	return &TokenService{secretKey: secretKey, tokenDuration: tokenDuration, emailVerificationToken: emailVerificationToken, forgotPasswordToken: forgotPasswordToken}
}

func (s *TokenService) GenerateToken(userID, role string, t time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"iat":    t.Unix(),
		"exp":    t.Add(s.tokenDuration).Unix(),
	})

	return token.SignedString([]byte(s.secretKey))
}

func (s *TokenService) ValidateToken(tokenString string) (auth.TokenPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return auth.TokenPayload{}, err
	}

	if !token.Valid {
		return auth.TokenPayload{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return auth.TokenPayload{}, fmt.Errorf("invalid token claims")
	}

	userIDClaim, exists := claims["userID"]
	if !exists || userIDClaim == nil {
		return auth.TokenPayload{}, fmt.Errorf("userID claim is missing or nil")
	}

	roleClaim, exists := claims["role"]
	if !exists || roleClaim == nil {
		return auth.TokenPayload{}, fmt.Errorf("role claim is missing or nil")
	}

	userID, ok := userIDClaim.(string)
	if !ok {
		return auth.TokenPayload{}, fmt.Errorf("invalid userID type")
	}

	role, ok := roleClaim.(string)
	if !ok {
		return auth.TokenPayload{}, fmt.Errorf("invalid role type")
	}

	return auth.TokenPayload{
		UserID: userID,
		Role:   role,
	}, nil
}

func (s *TokenService) GenerateTokenForForgotPassword(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(s.forgotPasswordToken).Unix(),
	})

	return token.SignedString([]byte(s.secretKey))
}

func (s *TokenService) ValidateForgotPasswordToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	emailClaim, exists := claims["email"]
	if !exists || emailClaim == nil {
		return "", fmt.Errorf("email claim is missing or nil")
	}

	email, ok := emailClaim.(string)
	if !ok {
		return "", fmt.Errorf("invalid email type")
	}

	return email, nil
}

func (s *TokenService) GenerateVerificationToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(s.tokenDuration).Unix(),
	})

	return token.SignedString([]byte(s.secretKey))
}

func (s *TokenService) ValidateVerifyEmailToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	emailClaim, exists := claims["email"]
	if !exists || emailClaim == nil {
		return "", fmt.Errorf("email claim is missing or nil")
	}

	email, ok := emailClaim.(string)
	if !ok {
		return "", fmt.Errorf("invalid email type")
	}

	return email, nil
}
