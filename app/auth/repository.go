package auth

import (
	"context"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
