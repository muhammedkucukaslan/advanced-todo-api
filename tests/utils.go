package test

import (
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func IsErrorStatusCode(code int) bool {
	return code >= 400 && code < 600
}

func VerifyErrorResponse(t *testing.T, body io.ReadCloser, expectedError error) {
	assert.NotEmpty(t, body, "response body should not be empty for error cases")

	var errResp domain.Error
	err := json.NewDecoder(body).Decode(&errResp)
	require.NoError(t, err, "failed to decode error response")

	assert.NotEmpty(t, errResp.Message, "error message should not be empty")
	assert.Equal(t, expectedError.Error(), errResp.Message, "error message should match")
}

func CreateTestContainer(t *testing.T, ctx context.Context) (*postgres.PostgresContainer, string) {
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
