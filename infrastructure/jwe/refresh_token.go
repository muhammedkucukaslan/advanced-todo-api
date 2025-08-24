package jwe

import (
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

func (s *Service) GenerateAuthRefreshToken(userID, role string) (string, error) {
	claims := AuthClaims{
		UserID: userID,
		Role:   role,
		Iat:    time.Now().Unix(),
		Exp:    time.Now().Add(s.authRefreshTokenDuration).Unix(),
	}
	return encryptClaims(s.refreshTokenEncrypter, claims)
}

func (s *Service) ValidateAuthRefreshToken(tokenString string) (*auth.TokenPayload, error) {
	var claims AuthClaims
	if err := decryptClaims(tokenString, s.refreshTokenEncryptionKey, &claims); err != nil {
		return nil, err
	}
	if isExpired(claims.Exp) {
		return nil, domain.ErrExpiredToken
	}

	return &auth.TokenPayload{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
