package testuser

import (
	"context"
	"net/http"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/user"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	testUtils "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCurrentUserHandler(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	postgresContainer, connStr := testUtils.CreatePostgresTestContainer(t, ctx)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate postgres container")
	}()

	repo, err := postgresRepo.NewRepository(connStr)
	require.NoError(t, err, "failed to create repository")
	runMigrations(t, connStr)
	setupTestUser(t, connStr)

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
			"valid user ID request",
			args{
				ctx: validCtx,
				req: validGetCurrentUserRequest,
			},
			mockUser,
			nil,
			http.StatusOK,
		},
		{
			"invalid user ID request",
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
		Id:              domain.RealUserId,
		FullName:        domain.TestUser.FullName,
		Email:           domain.TestUser.Email,
		Role:            domain.TestUser.Role,
		IsEmailVerified: domain.TestUser.IsEmailVerified,
	}

	validCtx := context.WithValue(context.Background(), domain.UserIDKey, domain.RealUserId)
	invalidCtx := context.WithValue(context.Background(), domain.UserIDKey, domain.FakeUserId)

	return mockUser, validCtx, invalidCtx
}
