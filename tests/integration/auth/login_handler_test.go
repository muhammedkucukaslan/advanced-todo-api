package testauth

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/validator"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {

	repo, err := postgres.NewRepository(os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	tokenService := jwt.NewTokenService(os.Getenv("JWT_SECRET_KEY"), time.Hour*24, time.Minute*10, time.Minute*10)
	logger := slog.NewLogger()
	validator := validator.NewValidator(logger)

	handler := auth.NewLoginHandler(repo, tokenService, validator, logger)

	type args struct {
		ctx context.Context
		req *auth.LoginRequest
	}

	userPassword := os.Getenv("TEST_USER_PASSWORD")

	tests := []struct {
		name    string
		args    args
		want    *auth.LoginResponse
		code    int
		wantErr error
	}{
		{"valid login request",
			args{context.Background(), &auth.LoginRequest{
				Email:    domain.TestUser.Email,
				Password: userPassword,
			}},
			&auth.LoginResponse{
				Token: "valid-token",
			},
			http.StatusOK,
			nil,
		},
		{"invalid email request",
			args{context.Background(), &auth.LoginRequest{
				Email:    "invalid-email",
				Password: userPassword,
			}},
			nil,
			http.StatusBadRequest,
			domain.ErrInvalidRequest,
		},
		{"invalid password request",
			args{context.Background(), &auth.LoginRequest{
				Email:    domain.TestUser.Email,
				Password: "wrong-password",
			}},
			nil,
			http.StatusBadRequest,
			domain.ErrInvalidCredentials,
		},
		{"user not found request",
			args{context.Background(), &auth.LoginRequest{
				Email:    "notfound@example.com",
				Password: "any-password",
			}},
			nil,
			http.StatusNotFound,
			domain.ErrEmailNotFound,
		},
		{"empty email request",
			args{context.Background(), &auth.LoginRequest{
				Email:    "",
				Password: userPassword,
			}},
			nil,
			http.StatusBadRequest,
			domain.ErrInvalidRequest,
		},
		{"empty password request",
			args{context.Background(), &auth.LoginRequest{
				Email:    domain.TestUser.Email,
				Password: "",
			}},
			nil,
			http.StatusBadRequest,
			domain.ErrInvalidRequest,
		},
		{"too short password request",
			args{context.Background(), &auth.LoginRequest{
				Email:    domain.TestUser.Email,
				Password: "short",
			}},
			nil,
			http.StatusBadRequest,
			domain.ErrInvalidRequest,
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
