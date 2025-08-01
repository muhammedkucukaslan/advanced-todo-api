package testtodo

import (
	"context"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type MockRepository struct {
}

func (m *MockRepository) CreateTodo(ctx context.Context, todo *domain.Todo) error {

	return nil
}

func (m *MockRepository) UpdateTodo(ctx context.Context, id uuid.UUID, title string) error {
	if id == uuid.Nil || title == "" {
		return domain.ErrInvalidRequest
	}
	return nil
}
