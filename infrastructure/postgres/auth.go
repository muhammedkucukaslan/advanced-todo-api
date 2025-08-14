package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (id, fullname, email, password, role) VALUES ($1, $2, $3, $4, $5)",
		user.Id, user.FullName, user.Email, user.Password, user.Role)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // Unique violation
			return domain.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, email, role, password FROM users WHERE email = $1", email)

	if err := row.Err(); err != nil {
		return nil, err
	}

	var user domain.User
	err := row.Scan(&user.Id, &user.Email, &user.Role, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEmailNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) SaveRefreshToken(ctx context.Context, record *domain.RefreshToken) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)",
		record.Id, record.UserID, record.Token, record.ExpiresAt, record.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpsertRefreshToken(ctx context.Context, record *domain.RefreshToken) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id) DO UPDATE SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at, created_at = EXCLUDED.created_at",
		record.Id, record.UserID, record.Token, record.ExpiresAt, record.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE token = $1", token)
	if err != nil {
		return err
	}
	return nil
}
