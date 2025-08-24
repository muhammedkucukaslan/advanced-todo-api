package jwe

import (
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

func (s *Service) GenerateAuthAccessToken(userID, role string) (string, error) {
	claims := AuthClaims{
		UserID: userID,
		Role:   role,
		Iat:    time.Now().Unix(),
		Exp:    time.Now().Add(s.authAccessTokenDuration).Unix(),
	}
	return encryptClaims(s.accessTokenEncrypter, claims)
}

func (s *Service) ValidateAuthAccessToken(tokenString string) (*auth.TokenPayload, error) {
	var claims AuthClaims
	if err := decryptClaims(tokenString, s.accessTokenEncryptionKey, &claims); err != nil {
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
