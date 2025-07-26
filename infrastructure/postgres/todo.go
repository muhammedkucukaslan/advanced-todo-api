package postgres

import (
	"context"

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
