package postgres

import (
	"context"
	"time"

	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

func (r *Repository) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO todos (user_id, id, title, completed)
		VALUES ($1, $2, $3, $4)
	`, todo.UserId, todo.Id, todo.Title, todo.Completed)
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

func (r *Repository) GetById(ctx context.Context, id uuid.UUID) (*todo.GetTodoByIdResponse, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, title, completed, created_at, completed_at
		FROM todos
		WHERE id = $1
	`, id)

	var resp todo.GetTodoByIdResponse
	var completedAt sql.NullTime
	if err := row.Scan(&resp.Id, &resp.Title, &resp.Completed, &resp.CreatedAt, &completedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrTodoNotFound
		}
		return nil, err
	}
	if completedAt.Valid {
		resp.CompletedAt = completedAt.Time
	} else {
		resp.CompletedAt = time.Time{}
	}
	return &resp, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM todos WHERE id = $1`, id)
	return err
}

func (r *Repository) GetTodosByUserID(ctx context.Context, userID uuid.UUID) (*todo.GetTodosResponse, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, completed, created_at, completed_at
		FROM todos
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos todo.GetTodosResponse
	for rows.Next() {
		var resp domain.Todo
		var completedAt sql.NullTime
		if err := rows.Scan(&resp.Id, &resp.Title, &resp.Completed, &resp.CreatedAt, &completedAt); err != nil {
			return nil, err
		}
		if completedAt.Valid {
			resp.CompletedAt = completedAt.Time
		} else {
			resp.CompletedAt = time.Time{}
		}
		todos = append(todos, resp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &todos, nil
}
