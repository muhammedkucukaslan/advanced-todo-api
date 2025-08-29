package jwe

import (
	"encoding/json"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"gopkg.in/square/go-jose.v2"
)

func is32ByteKey(key string) bool {
	return len([]byte(key)) == 32
}

func encryptClaims[T any](encrypter jose.Encrypter, claims T) (string, error) {
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	jwe, err := encrypter.Encrypt(payload)
	if err != nil {
		return "", err
	}
	return jwe.CompactSerialize()
}

func decryptClaims[T any](encryptedToken string, key []byte, claims *T) error {
	jwe, err := jose.ParseEncrypted(encryptedToken)
	if err != nil {
		return domain.ErrInvalidToken
	}
	decrypted, err := jwe.Decrypt(key)
	if err != nil {
		return domain.ErrInvalidToken
	}
	if err := json.Unmarshal(decrypted, claims); err != nil {
		return domain.ErrInvalidToken
	}
	return nil
}

func isExpired(exp int64) bool {
	return time.Now().Unix() > exp
}
