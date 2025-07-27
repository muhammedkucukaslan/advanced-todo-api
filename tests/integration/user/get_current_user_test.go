package testuser

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/user"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentUserHandler(t *testing.T) {

	repo, err := postgres.NewRepository(os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	getCurrentUserHandler := user.NewGetCurrentUserHandler(repo)

	mockUser, validCtx, invalidCtx := setupGetCurrentUserTestData()

	validGetCurrentUserRequest := &user.GetCurrentUserRequest{}

	type args struct {
		ctx context.Context
		req *user.GetCurrentUserRequest
	}

	tests := []struct {
		name    string
		args    args
		want    *user.GetCurrentUserResponse
		wantErr error
		code    int
	}{
		{
			"should get current user successfully",
			args{
				ctx: validCtx,
				req: validGetCurrentUserRequest,
			},
			mockUser,
			nil,
			http.StatusOK,
		},
		{
			"should return 404 for fake user ID",
			args{
				ctx: invalidCtx,
				req: validGetCurrentUserRequest,
			},
			nil,
			domain.ErrUserNotFound,
			http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, code, err := getCurrentUserHandler.Handle(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.code, code)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {

				assert.ErrorIs(t, err, tt.wantErr)

				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.FullName, got.FullName)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Role, got.Role)
				assert.Equal(t, tt.want.IsEmailVerified, got.IsEmailVerified)
				// i could use assert.Equal(t, tt.want, got).However, i got some problem on createdAt field.
			}
		})
	}
}

func setupGetCurrentUserTestData() (*user.GetCurrentUserResponse, context.Context, context.Context) {

	mockUser := &user.GetCurrentUserResponse{
		Id:              domain.TestUser.Id.String(),
		FullName:        domain.TestUser.FullName,
		Email:           domain.TestUser.Email,
		Role:            domain.TestUser.Role,
		IsEmailVerified: domain.TestUser.IsEmailVerified,
	}

	validCtx := context.WithValue(context.Background(), domain.UserIDKey, domain.RealUserId)
	invalidCtx := context.WithValue(context.Background(), domain.UserIDKey, domain.FakeUserId)

	return mockUser, validCtx, invalidCtx
}
