package jwe

import (
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

// These functions are used for "forgot password" and "email verification"

func (s *Service) GenerateSecureEmailToken(email string) (string, error) {
	claims := EmailClaims{
		Email: email,
		Iat:   time.Now().Unix(),
		Exp:   time.Now().Add(s.secureEmailTokenDuration).Unix(),
	}
	return encryptClaims(s.secureEmailEncrypter, claims)
}

func (s *Service) ValidateSecureEmailToken(tokenString string) (string, error) {
	var claims EmailClaims
	if err := decryptClaims(tokenString, s.secureEmailEncryptionKey, &claims); err != nil {
		return "", err
	}
	if isExpired(claims.Exp) {
		return "", domain.ErrExpiredToken
	}
	return claims.Email, nil
}
