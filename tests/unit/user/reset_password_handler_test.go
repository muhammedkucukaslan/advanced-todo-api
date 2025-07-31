package testuser

import (
	"context"
	"net/http"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/user"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	mock "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/assert"
)

func TestResetPasswordHandler(t *testing.T) {

	handler := user.NewResetPasswordHandler(NewMockRepository(), mock.NewMockTokenService(), mock.NewMockLogger(), mock.NewMockValidator())

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
		{"valid request", args{
			ctx: context.Background(),
			req: &user.ResetPasswordRequest{
				Token:    "validToken",
				Password: "validPassword123",
			},
		}, http.StatusNoContent, nil},
		{"too short password", args{
			ctx: context.Background(),
			req: &user.ResetPasswordRequest{
				Token:    "validToken",
				Password: "short",
			},
		}, http.StatusBadRequest, domain.ErrPasswordTooShort},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, code, err := handler.Handle(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.code, code)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
