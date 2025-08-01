package testtodo

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	fmt.Println("Running userintegration tests...")
	err := godotenv.Load("../../../.env")
	if err != nil {
		panic("Failed to load .env file: " + err.Error())
	}

	code := m.Run()

	os.Exit(code)
}

func createTestContainer(t *testing.T, ctx context.Context) (*postgres.PostgresContainer, string) {
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	require.NoError(t, err)
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	return postgresContainer, connStr
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
	`

	_, err = db.Exec(createTableQuery)
	require.NoError(t, err)
}
