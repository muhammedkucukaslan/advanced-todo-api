package unittest_user

import (
	"context"
	"net/http"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/user"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	mock "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/assert"
)

func TestChangePasswordHandler(t *testing.T) {

	handler := user.NewChangePasswordHandler(NewMockRepository(), mock.NewMockValidator())

	ctx := context.WithValue(context.Background(), domain.UserIDKey, domain.TestUser.Id.String())

	type args struct {
		ctx context.Context
		req *user.ChangePasswordRequest
	}

	tests := []struct {
		name    string
		args    args
		code    int
		wantErr error
	}{
		{"valid request", args{
			ctx: ctx,
			req: &user.ChangePasswordRequest{
				OldPassword: "password123",
				NewPassword: "newPassword123",
			},
		}, http.StatusNoContent, nil},
		{"wrong old password", args{
			ctx: ctx,
			req: &user.ChangePasswordRequest{
				OldPassword: "wrongOldPassword",
				NewPassword: "newPassword123",
			},
		}, http.StatusBadRequest, domain.ErrInvalidCredentials},
		{"too short old password", args{
			ctx: ctx,
			req: &user.ChangePasswordRequest{
				OldPassword: "short",
				NewPassword: "newPassword123",
			},
		}, http.StatusBadRequest, domain.ErrPasswordTooShort},
		{"too short new password", args{
			ctx: ctx,
			req: &user.ChangePasswordRequest{
				OldPassword: "password123",
				NewPassword: "short",
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
