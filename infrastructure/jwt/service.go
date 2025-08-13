package jwt

import (
	"errors"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	SecretKey                 string
	AuthAccessTokenDuration   time.Duration
	EmailVerificationDuration time.Duration
	ForgotPasswordDuration    time.Duration
}

type Service struct {
	secretKey                 []byte
	AuthAccessTokenDuration   time.Duration
	emailVerificationDuration time.Duration
	forgotPasswordDuration    time.Duration
}

func NewJWTTokenService(config Config) *Service {
	return &Service{
		secretKey:                 []byte(config.SecretKey),
		AuthAccessTokenDuration:   config.AuthAccessTokenDuration,
		emailVerificationDuration: config.EmailVerificationDuration,
		forgotPasswordDuration:    config.ForgotPasswordDuration,
	}
}

func (s *Service) GenerateAuthAccessToken(userID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(s.AuthAccessTokenDuration).Unix(),
	})

	return token.SignedString(s.secretKey)
}

func (s *Service) ValidateAuthAccessToken(tokenString string) (*auth.TokenPayload, error) {
	token, err := jwt.Parse(tokenString, s.keyFunc)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, domain.ErrInvalidTokenSignature
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, domain.ErrInvalidToken
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return s.validateAuthClaims(token)
}

func (s *Service) GenerateTokenForForgotPassword(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(s.forgotPasswordDuration).Unix(),
	})

	return token.SignedString(s.secretKey)
}

func (s *Service) ValidateForgotPasswordToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, s.keyFunc)

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}
	return s.validateEmailClaims(token)
}

func (s *Service) GenerateEmailVerificationToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(s.emailVerificationDuration).Unix(),
	})

	return token.SignedString(s.secretKey)
}

func (s *Service) ValidateVerifyEmailToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, s.keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", domain.ErrExpiredToken
		}
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return s.validateEmailClaims(token)
}

func (s *Service) keyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}

	return s.secretKey, nil
}

func (s *Service) validateAuthClaims(token *jwt.Token) (*auth.TokenPayload, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDClaim, exists := claims["userID"]
	if !exists || userIDClaim == nil {
		return nil, errors.New("userID claim is missing or nil")
	}

	roleClaim, exists := claims["role"]
	if !exists || roleClaim == nil {
		return nil, errors.New("role claim is missing or nil")
	}

	userID, ok := userIDClaim.(string)
	if !ok {
		return nil, errors.New("invalid userID type")
	}

	role, ok := roleClaim.(string)
	if !ok {
		return nil, errors.New("invalid role type")
	}

	return &auth.TokenPayload{
		UserID: userID,
		Role:   role,
	}, nil
}

func (s *Service) validateEmailClaims(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	emailClaim, exists := claims["email"]
	if !exists || emailClaim == nil {
		return "", errors.New("email claim is missing or nil")
	}

	email, ok := emailClaim.(string)
	if !ok {
		return "", errors.New("invalid email type")
	}

	return email, nil
}
