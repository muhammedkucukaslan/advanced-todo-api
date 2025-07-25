package auth

import (
	"context"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
