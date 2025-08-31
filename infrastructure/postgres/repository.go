package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(databaseUrl string) *Repository {
	if databaseUrl == "" {
		panic("Database URL is empty")
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil
	}

	if err := db.Ping(); err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	runTableMigrations(db)

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

		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id              UUID PRIMARY KEY,
			user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token           TEXT NOT NULL UNIQUE,
			expires_at      TIMESTAMP NOT NULL,
			created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(user_id)
		);
	`
	_, err := db.Exec(createTableQuery)

	if err != nil {
		panic("Failed to create tables: " + err.Error())
	}

}
