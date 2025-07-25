package domain

import (
	"context"

	"github.com/google/uuid"
)

type ContextKey string

const (
	UserIDKey ContextKey = "userID"
	RoleKey   ContextKey = "role"
	TokenKey  ContextKey = "token"
)

func GetUserID(ctx context.Context) uuid.UUID {
	userID := ctx.Value(UserIDKey).(string)
	return uuid.MustParse(userID)
}

func GetRole(ctx context.Context) string {
	return ctx.Value(RoleKey).(string)
}
