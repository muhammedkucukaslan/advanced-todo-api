package testuser

import (
	"context"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/user"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	mock "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/assert"
)

func TestResetPasswordHandler_Handle(t *testing.T) {

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
		}, 204, nil},
		{"too short password", args{
			ctx: context.Background(),
			req: &user.ResetPasswordRequest{
				Token:    "validToken",
				Password: "short",
			},
		}, 400, domain.ErrPasswordTooShort},
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
