package postgres

import (
	"database/sql"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(databaseUrl string) (*Repository, error) {
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if os.Getenv("ENV") != "production" {
		if err := runMigrations(db); err != nil {
			return nil, err
		}
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func rollbackTx(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		log.Printf("failed to rollback transaction: %v", err)
	}
}

func runMigrations(db *sql.DB) error {

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			fullname VARCHAR(255),
			role VARCHAR(10) NOT NULL CHECK (role IN ('USER', 'ADMIN')),
			password VARCHAR(200) NOT NULL,
			email VARCHAR(200) NOT NULL UNIQUE,
			is_email_verified BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);


		CREATE TABLE IF NOT EXISTS todos (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP DEFAULT NULL
		);
	`
	if _, err := db.Exec(createTableQuery); err != nil {
		return err
	}
	hashedPassword, err := domain.HashPassword(domain.TestUser.Password)
	if err != nil {
		return err
	}

	addUserQuery := `INSERT INTO users (id, fullname, email, password, role) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (email) DO NOTHING`
	_, err = db.Exec(addUserQuery,
		domain.TestUser.Id,
		domain.TestUser.FullName,
		domain.TestUser.Email,
		hashedPassword,
		domain.TestUser.Role,
	)

	adminHashedPassword, err := domain.HashPassword("admin123")
	if err != nil {
		return err
	}
	addAdminQuery := `INSERT INTO users (id, fullname, email, password, role) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (email) DO NOTHING`
	_, err = db.Exec(addAdminQuery,
		uuid.New(),
		"ADMIN",
		"admin@admin.com",
		adminHashedPassword,
		"ADMIN",
	)

	return err

}
