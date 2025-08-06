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
	mock "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	testUtils "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForgotPasswordHandler(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	postgresContainer, connStr := testUtils.CreateTestContainer(t, ctx)
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
	mockEmailService := &mock.MockEmailService{}

	handler := user.NewForgotPasswordHandler(repo, mockEmailService, tokenService, logger, validator)

	type args struct {
		ctx context.Context
		req *user.ForgotPasswordRequest
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
				req: &user.ForgotPasswordRequest{
					Email: domain.TestUser.Email,
				},
			},
			code:    http.StatusNoContent,
			wantErr: nil,
		},
		{
			name: "invalid email",
			args: args{
				ctx: ctx,
				req: &user.ForgotPasswordRequest{
					Email: "invalid-email",
				},
			},
			code:    http.StatusBadRequest,
			wantErr: domain.ErrInvalidRequest,
		},
		{
			name: "email not found",
			args: args{
				ctx: ctx,
				req: &user.ForgotPasswordRequest{
					Email: "notfound@example.com",
				},
			},
			// If the email does not exist, we still return a success response
			// to prevent email enumeration attacks.
			code:    http.StatusNoContent,
			wantErr: nil,
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
