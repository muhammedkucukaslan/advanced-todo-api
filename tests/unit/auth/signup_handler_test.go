package testauth

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	mock "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/assert"
)

func Test_Name(t *testing.T) {

	handler := auth.NewSignupHandler(&auth.SignupConfig{
		RefreshTokenCookieDuration: time.Hour * 24 * 30,
		Secure:                     false,
		Repo:                       NewMockRepository(),
		TokenService:               mock.NewMockTokenService(),
		CookieService:              mock.NewMockCookieService(),
		EmailService:               mock.NewMockEmailService(),
		Validator:                  mock.NewMockValidator(),
		Logger:                     mock.NewMockLogger(),
	})

	type args struct {
		ctx context.Context
		req *auth.SignupRequest
	}

	tests := []struct {
		name    string
		args    args
		want    *auth.SignupResponse
		code    int
		wantErr error
	}{
		{"valid signup", args{context.Background(), &auth.SignupRequest{
			FullName: domain.TestUser.FullName,
			Password: "validpassword",
			Email:    domain.TestUser.Email,
		}}, &auth.SignupResponse{
			AccessToken: domain.MockToken,
		}, http.StatusCreated, nil},
		{"too short fullName", args{context.Background(), &auth.SignupRequest{
			FullName: "sh",
			Password: "validpassword",
			Email:    domain.TestUser.Email,
		}}, nil, http.StatusBadRequest, domain.ErrTooShortFullName},
		{"too short password", args{context.Background(), &auth.SignupRequest{
			FullName: domain.TestUser.FullName,
			Password: "short",
			Email:    domain.TestUser.Email,
		}}, nil, http.StatusBadRequest, domain.ErrPasswordTooShort},
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
				assert.Equal(t, tt.want, got)
			}

		})
	}
}
