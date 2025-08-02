package todo

import (
	"context"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, todo *domain.Todo) error
	UpdateTodo(ctx context.Context, id uuid.UUID, title string) error
	GetById(ctx context.Context, id uuid.UUID) (*GetTodoByIdResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetTodosByUserID(ctx context.Context, userID uuid.UUID) (*GetTodosResponse, error)
}
