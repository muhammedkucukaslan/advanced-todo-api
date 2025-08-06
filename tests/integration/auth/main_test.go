package testauth

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	fmt.Println("Running auth integration tests...")

	code := m.Run()

	os.Exit(code)
}

func setupTestUser(t *testing.T, connStr string) {

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	hashedPassword, err := domain.HashPassword(domain.TestUser.Password)
	require.NoError(t, err)

	query := `INSERT INTO users (id, fullname, email, password, role) VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (email) DO NOTHING`
	_, err = db.Exec(query,
		domain.TestUser.Id,
		domain.TestUser.FullName,
		domain.TestUser.Email,
		hashedPassword,
		domain.TestUser.Role,
	)
	require.NoError(t, err)

}

func runMigrations(t *testing.T, connStr string) {
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

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

		TRUNCATE TABLE users CASCADE;
	`

	_, err = db.Exec(createTableQuery)
	require.NoError(t, err)
}
