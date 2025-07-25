package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

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
