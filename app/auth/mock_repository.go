package auth

import (
	"context"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type AuthMockRepository struct {
	// Add fields to simulate the behavior of the real repository
}

func (r *AuthMockRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Simulate the behavior of getting a user by email
	return &domain.User{Email: email}, nil // or return an error if needed
}

func NewAuthMockRepository() *AuthMockRepository {
	return &AuthMockRepository{}
}
