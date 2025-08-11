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
	jwtInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	postgresInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	slogInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/validator"
	mock "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	testUtils "github.com/muhammedkucukaslan/advanced-todo-api/tests"
)

func TestSignupHandler(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	postgresContainer, connStr := testUtils.CreatePostgresTestContainer(t, ctx)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate postgres container")
	}()

	repo := postgresInfra.NewRepository(connStr)
	runMigrations(t, connStr)
	setupTestUser(t, connStr)

	tokenService := jwtInfra.NewTokenService("test-secret-key", time.Hour*24, time.Minute*10, time.Minute*10)
	logger := slogInfra.NewLogger()
	validator := validator.NewValidator(logger)
	mailService := mock.NewMockEmailService()

	handler := auth.NewSignupHandler(repo, tokenService, mailService, validator, logger)

	type args struct {
		ctx context.Context
		req *auth.SignupRequest
	}

	mockUser := &domain.User{
		FullName: "Different User Nick Name",
		Email:    "different@example.com",
		Password: "different-password",
	}

	tests := []struct {
		name    string
		args    args
		want    *auth.SignupResponse
		code    int
		wantErr error
	}{
		{
			name: "valid signup",
			args: args{
				ctx: context.Background(),
				req: &auth.SignupRequest{
					FullName: mockUser.FullName,
					Email:    mockUser.Email,
					Password: mockUser.Password,
				},
			},
			want: &auth.SignupResponse{
				Token: "valid-token",
			},
			code:    http.StatusCreated,
			wantErr: nil,
		},
		{
			name: "invalid email ",
			args: args{
				ctx: context.Background(),
				req: &auth.SignupRequest{
					FullName: mockUser.FullName,
					Email:    "invalid-email",
					Password: mockUser.Password,
				},
			},
			want:    nil,
			code:    http.StatusBadRequest,
			wantErr: domain.ErrInvalidRequest,
		},
		{
			name: "empty email ",
			args: args{
				ctx: context.Background(),
				req: &auth.SignupRequest{
					FullName: mockUser.FullName,
					Password: mockUser.Password,
				},
			},
			want:    nil,
			code:    http.StatusBadRequest,
			wantErr: domain.ErrInvalidRequest,
		},
		{
			name: "duplicate email ",
			args: args{
				ctx: context.Background(),
				req: &auth.SignupRequest{
					FullName: mockUser.FullName,
					Email:    domain.TestUser.Email,
					Password: mockUser.Password,
				},
			},
			want:    nil,
			code:    http.StatusConflict,
			wantErr: domain.ErrEmailAlreadyExists,
		},
		{
			name: "empty password ",
			args: args{
				ctx: context.Background(),
				req: &auth.SignupRequest{
					FullName: mockUser.FullName,
					Email:    mockUser.Email,
				},
			},
			want:    nil,
			code:    http.StatusBadRequest,
			wantErr: domain.ErrInvalidRequest,
		},
		{
			name: "empty full name ",
			args: args{
				ctx: context.Background(),
				req: &auth.SignupRequest{
					Email:    mockUser.Email,
					Password: mockUser.Password,
				},
			},
			want:    nil,
			code:    http.StatusBadRequest,
			wantErr: domain.ErrInvalidRequest,
		},
		{
			name: "too short full name ",
			args: args{
				ctx: context.Background(),
				req: &auth.SignupRequest{
					FullName: "ab",
					Email:    mockUser.Email,
					Password: mockUser.Password,
				},
			},
			want:    nil,
			code:    http.StatusBadRequest,
			wantErr: domain.ErrTooShortFullName,
		},
		{
			name: "too short password ",
			args: args{
				ctx: context.Background(),
				req: &auth.SignupRequest{
					FullName: mockUser.FullName,
					Email:    mockUser.Email,
					Password: "short",
				},
			},
			want:    nil,
			code:    http.StatusBadRequest,
			wantErr: domain.ErrPasswordTooShort,
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
				assert.Equal(t, domain.TestUser.Role, payload.Role)
			}
		})
	}
}
