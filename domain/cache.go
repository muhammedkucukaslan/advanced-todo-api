package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

func NewTodoCacheKey(userId uuid.UUID) string {
	return "todos:" + userId.String()
}
