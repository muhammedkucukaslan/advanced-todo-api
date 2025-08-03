package testuser

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/user"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResetPasswordHandler(t *testing.T) {
	ctx := context.Background()

	postgresContainer, connStr := createTestContainer(t, ctx)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate postgres container")
	}()

	repo, err := postgresRepo.NewRepository(connStr)
	require.NoError(t, err, "failed to create repository")
	runMigrations(t, connStr)
	setupTestUser(t, connStr)

	tokenService := jwt.NewTokenService(domain.MockJWTTestKey, time.Hour*24, time.Minute*10, time.Minute*10)
	logger := slog.NewLogger()
	validator := validator.NewValidator(logger)

	mockJWTToken, err := tokenService.GenerateTokenForForgotPassword(domain.TestUser.Email)
	require.NoError(t, err, "failed to generate mock JWT token")

	handler := user.NewResetPasswordHandler(repo, tokenService, logger, validator)

	type args struct {
		ctx context.Context
		req *user.ResetPasswordRequest
	}

	tests := []struct {
		name    string
		args    args
		code    int
		wantErr error
	}{
		{
			name: "valid request",
			args: args{
				ctx: ctx,
				req: &user.ResetPasswordRequest{
					Token:    mockJWTToken,
					Password: "new-password",
				},
			},
			code:    http.StatusNoContent,
			wantErr: nil,
		},
		{
			name: "invalid token",
			args: args{
				ctx: ctx,
				req: &user.ResetPasswordRequest{
					Token:    "invalid-token",
					Password: "new-password",
				},
			},
			code:    http.StatusUnauthorized,
			wantErr: domain.ErrUnauthorized,
		},
		{
			name: "too short password",
			args: args{
				ctx: ctx,
				req: &user.ResetPasswordRequest{
					Token:    mockJWTToken,
					Password: "123",
				},
			},
			code:    http.StatusBadRequest,
			wantErr: domain.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, code, err := handler.Handle(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.code, code)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
