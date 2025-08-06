package testauth

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	postgresInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/validator"
	testUtils "github.com/muhammedkucukaslan/advanced-todo-api/tests"
)

func TestLoginHandler(t *testing.T) {
	ctx := context.Background()

	postgresContainer, connStr := testUtils.CreateTestContainer(t, ctx)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate postgres container")
	}()

	repo, err := postgresInfra.NewRepository(connStr)
	require.NoError(t, err, "failed to create repository")
	runMigrations(t, connStr)
	setupTestUser(t, connStr)

	tokenService := jwt.NewTokenService("test-secret-key", time.Hour*24, time.Minute*10, time.Minute*10)
	logger := slog.NewLogger()
	validator := validator.NewValidator(logger)

	handler := auth.NewLoginHandler(repo, tokenService, validator, logger)

	type args struct {
		ctx context.Context
		req *auth.LoginRequest
	}

	tests := []struct {
		name    string
		args    args
		want    *auth.LoginResponse
		code    int
		wantErr error
	}{
		{
			name: "valid login request",
			args: args{
				ctx: context.Background(),
				req: &auth.LoginRequest{
					Email:    domain.TestUser.Email,
					Password: domain.TestUser.Password,
				},
			},
			want: &auth.LoginResponse{
				Token: "valid-token",
			},
			code:    http.StatusOK,
			wantErr: nil,
		},
		{
			name: "invalid email request",
			args: args{
				ctx: context.Background(),
				req: &auth.LoginRequest{
					Email:    "invalid-email",
					Password: domain.TestUser.Password,
				},
			},
			want:    nil,
			code:    http.StatusBadRequest,
			wantErr: domain.ErrInvalidRequest,
		},
		{
			name: "invalid password request",
			args: args{
				ctx: context.Background(),
				req: &auth.LoginRequest{
					Email:    domain.TestUser.Email,
					Password: "wrong-password",
				},
			},
			want:    nil,
			code:    http.StatusBadRequest,
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "user not found request",
			args: args{
				ctx: context.Background(),
				req: &auth.LoginRequest{
					Email:    "notfound@example.com",
					Password: "any-password",
				},
			},
			want:    nil,
			code:    http.StatusNotFound,
			wantErr: domain.ErrEmailNotFound,
		},
		{
			name: "empty email request",
			args: args{
				ctx: context.Background(),
				req: &auth.LoginRequest{
					Email:    "",
					Password: domain.TestUser.Password,
				},
			},
			want:    nil,
			code:    http.StatusBadRequest,
			wantErr: domain.ErrInvalidRequest,
		},
		{
			name: "empty password request",
			args: args{
				ctx: context.Background(),
				req: &auth.LoginRequest{
					Email:    domain.TestUser.Email,
					Password: "",
				},
			},
			want:    nil,
			code:    http.StatusBadRequest,
			wantErr: domain.ErrInvalidRequest,
		},
		{
			name: "too short password request",
			args: args{
				ctx: context.Background(),
				req: &auth.LoginRequest{
					Email:    domain.TestUser.Email,
					Password: "short",
				},
			},
			want:    nil,
			code:    http.StatusBadRequest,
			wantErr: domain.ErrInvalidRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, code, err := handler.Handle(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.code, code)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.NotEmpty(t, got.Token)
				payload, err := tokenService.ValidateToken(got.Token)
				assert.NoError(t, err)
				assert.NotNil(t, payload)
				assert.Equal(t, domain.RealUserId, payload.UserID)
				assert.Equal(t, domain.TestUser.Role, payload.Role)
			}
		})
	}
}
