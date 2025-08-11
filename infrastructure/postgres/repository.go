package postgres

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(databaseUrl string) *Repository {
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil
	}

	if err := db.Ping(); err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	runTableMigrations(db)

	if !domain.IsProdEnv() {
		runTestUserMigrations(db)
	}

	return &Repository{db: db}
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func rollbackTx(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		log.Printf("failed to rollback transaction: %v", err)
	}
}

func runTableMigrations(db *sql.DB) {
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
	_, err := db.Exec(createTableQuery)

	if err != nil {
		panic("Failed to create tables: " + err.Error())
	}

}

func runTestUserMigrations(db *sql.DB) {
	tx, err := db.Begin()

	if err != nil {
		panic("Failed to begin transaction: " + err.Error())
	}
	defer rollbackTx(tx)

	userHashedPassword, err := domain.HashPassword(domain.TestUser.Password)
	if err != nil {
		panic("Failed to hash password: " + err.Error())
	}

	addUserQuery := `INSERT INTO users (id, fullname, email, password, role) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (email) DO NOTHING`
	_, err = tx.Exec(addUserQuery,
		domain.TestUser.Id,
		domain.TestUser.FullName,
		domain.TestUser.Email,
		userHashedPassword,
		domain.TestUser.Role,
	)

	adminHashedPassword, err := domain.HashPassword("admin123")
	if err != nil {
		panic("Failed to hash password: " + err.Error())
	}
	addAdminQuery := `INSERT INTO users (id, fullname, email, password, role) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (email) DO NOTHING`
	_, err = tx.Exec(addAdminQuery,
		uuid.New(),
		"Admin User",
		"admin@admin.com",
		adminHashedPassword,
		"ADMIN",
	)
	if err != nil {
		panic("Failed to insert admin user: " + err.Error())
	}
	if tx.Commit() != nil {
		panic("Failed to commit transaction: " + tx.Commit().Error())
	}
}
