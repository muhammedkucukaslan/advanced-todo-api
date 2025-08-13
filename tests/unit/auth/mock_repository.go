package testauth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type MockRepository struct {
}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

func (m *MockRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == "test@example.com" {
		return &domain.User{
			Id:       uuid.New(),
			Email:    "test@example.com",
			FullName: "Test User",
		}, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockRepository) SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	return nil
}
