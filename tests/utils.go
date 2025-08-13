package test

import (
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	jwtInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
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

func CreatePostgresTestContainer(t *testing.T, ctx context.Context) (*tcpostgres.PostgresContainer, string) {
	postgresContainer, err := tcpostgres.Run(ctx,
		"postgres:15",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
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

func CreateRedisTestContainer(t *testing.T, ctx context.Context) (*tcredis.RedisContainer, string) {
	redisContainer, err := tcredis.Run(ctx, "redis:7")

	require.NoError(t, err, "failed to start Redis container")

	uri, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err, "failed to get Redis connection string")

	return redisContainer, uri

}

func NewTestJWTTokenService() *jwtInfra.Service {
	return jwtInfra.NewJWTTokenService(jwtInfra.Config{
		AccessTokenSecretKey:      "test_secret_key",
		RefreshTokenSecretKey:     "test_refresh_secret_key",
		AuthAccessTokenDuration:   time.Hour * 24,
		AuthRefreshTokenDuration:  time.Hour * 24 * 30,
		EmailVerificationDuration: time.Minute * 10,
		ForgotPasswordDuration:    time.Minute * 10,
	})
}
