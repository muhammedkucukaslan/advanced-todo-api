package testtodo

import (
	"context"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
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

func (m *MockRepository) GetById(ctx context.Context, id uuid.UUID) (*todo.GetTodoByIdResponse, error) {
	return nil, nil
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
