package e2etest_todo

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	fmt.Println("Running tests in e2etest_todo package...")

	code := m.Run()

	os.Exit(code)
}

func setupTestUser(t *testing.T, connStr string) {
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	hashedPassword, err := domain.HashPassword(domain.TestUser.Password)
	require.NoError(t, err)

	query := "INSERT INTO users (id, fullname, email, password, role) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.Exec(query,
		domain.TestUser.Id,
		domain.TestUser.FullName,
		domain.TestUser.Email,
		hashedPassword,
		domain.TestUser.Role,
	)
	require.NoError(t, err)

}

func setupTestTodoWithTitle(t *testing.T, id uuid.UUID, connStr string) {
	t.Helper()

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO todos (id, user_id, title, completed)
		VALUES ($1, $2, $3, $4)
	`, id, domain.RealUserId, domain.TestTodo.Title, domain.TestTodo.Completed)
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

		CREATE TABLE IF NOT EXISTS todos (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP DEFAULT NULL
		);

        TRUNCATE TABLE todos CASCADE;
	`

	_, err = db.Exec(createTableQuery)
	require.NoError(t, err)
}
