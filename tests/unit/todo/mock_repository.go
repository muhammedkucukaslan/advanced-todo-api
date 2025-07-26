package testtodo

import (
	"context"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type MockRepository struct {
}

func (m *MockRepository) CreateTodo(ctx context.Context, todo *domain.Todo) error {

	return nil
}
