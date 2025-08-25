package jwe

import (
	"time"

	"gopkg.in/square/go-jose.v2"
)

type Config struct {
	AccessTokenEncryptionKey  string
	RefreshTokenEncryptionKey string
	SecureEmailEncryptionKey  string

	AuthAccessTokenDuration  time.Duration
	AuthRefreshTokenDuration time.Duration
	SecureEmailTokenDuration time.Duration
}

type Service struct {
	accessTokenEncryptionKey  []byte
	refreshTokenEncryptionKey []byte
	secureEmailEncryptionKey  []byte

	accessTokenEncrypter  jose.Encrypter
	refreshTokenEncrypter jose.Encrypter
	secureEmailEncrypter  jose.Encrypter

	authAccessTokenDuration  time.Duration
	authRefreshTokenDuration time.Duration
	secureEmailTokenDuration time.Duration
}

type AuthClaims struct {
	UserID string `json:"userID"`
	Role   string `json:"role"`
	Iat    int64  `json:"iat"`
	Exp    int64  `json:"exp"`
}

type EmailClaims struct {
	Email string `json:"email"`
	Iat   int64  `json:"iat"`
	Exp   int64  `json:"exp"`
}

func NewJWETokenService(config *Config) *Service {

	if !config.hasProperEncryptionKeys() {
		panic("encryption key must be 32 bytes for AES-256")
	}

	accessTokenEncrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{Algorithm: jose.DIRECT, Key: []byte(config.AccessTokenEncryptionKey)},
		(&jose.EncrypterOptions{}).
			WithType("JWE").
			WithContentType("JWT"),
	)

	if err != nil {
		panic(err)
	}

	refreshTokenEncrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{Algorithm: jose.DIRECT, Key: []byte(config.RefreshTokenEncryptionKey)},
		(&jose.EncrypterOptions{}).
			WithType("JWE").
			WithContentType("JWT"),
	)
	if err != nil {
		panic(err)
	}

	secureEmailEncrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{Algorithm: jose.DIRECT, Key: []byte(config.SecureEmailEncryptionKey)},
		(&jose.EncrypterOptions{}).
			WithType("JWE").
			WithContentType("JWT"),
	)
	if err != nil {
		panic(err)
	}

	return &Service{
		accessTokenEncryptionKey:  []byte(config.AccessTokenEncryptionKey),
		refreshTokenEncryptionKey: []byte(config.RefreshTokenEncryptionKey),
		secureEmailEncryptionKey:  []byte(config.SecureEmailEncryptionKey),

		secureEmailEncrypter:  secureEmailEncrypter,
		accessTokenEncrypter:  accessTokenEncrypter,
		refreshTokenEncrypter: refreshTokenEncrypter,

		authAccessTokenDuration:  config.AuthAccessTokenDuration,
		authRefreshTokenDuration: config.AuthRefreshTokenDuration,
		secureEmailTokenDuration: config.SecureEmailTokenDuration,
	}
}

func (c Config) hasProperEncryptionKeys() bool {
	return is32ByteKey(c.AccessTokenEncryptionKey) && is32ByteKey(c.RefreshTokenEncryptionKey) && is32ByteKey(c.SecureEmailEncryptionKey)
}
