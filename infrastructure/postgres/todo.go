package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

func (r *Repository) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO todos (user_id, id, title, completed, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, todo.UserId, todo.Id, todo.Title, todo.Completed, todo.CreatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *Repository) UpdateTodo(ctx context.Context, id uuid.UUID, title string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE todos
		SET title = $1
		WHERE id = $2
	`, title, id)
	return err
}
