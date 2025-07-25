package todo

import (
	"context"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, todo *domain.Todo) error
}
