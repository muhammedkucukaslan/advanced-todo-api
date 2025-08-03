package testtodo

import (
	"context"
	"net/http"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTodoHandler(t *testing.T) {
	ctx := context.Background()

	postgresContainer, connStr := createTestContainer(t, ctx)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate postgres container")
	}()

	repo, err := postgresRepo.NewRepository(connStr)
	require.NoError(t, err, "failed to create repository")
	runMigrations(t, connStr)
	setupTestUser(t, connStr)

	createTodoHandler := todo.NewCreateTodoHandler(repo)
	ctx = context.WithValue(context.Background(), domain.UserIDKey, domain.RealUserId)

	ctxWithFakeUserId := context.WithValue(context.Background(), domain.UserIDKey, domain.FakeUserId)

	type args struct {
		ctx context.Context
		req *todo.CreateTodoRequest
	}

	tests := []struct {
		name    string
		args    args
		code    int
		wantErr error
	}{
		{
			"valid creation",
			args{
				ctx: ctx,
				req: &todo.CreateTodoRequest{
					Title: "Test Todo",
				},
			},
			http.StatusCreated,
			nil,
		},
		{
			"invalid user ID request",
			args{
				ctx: ctxWithFakeUserId,
				req: &todo.CreateTodoRequest{
					Title: "Test Todo",
				},
			},
			http.StatusNotFound,
			domain.ErrUserNotFound,
		},
		{
			"invalid request",
			args{
				ctx: ctx,
				req: &todo.CreateTodoRequest{},
			},
			http.StatusBadRequest,
			domain.ErrInvalidRequest,
		},
		{
			"too short title",
			args{
				ctx: ctx,
				req: &todo.CreateTodoRequest{
					Title: "ab",
				},
			},
			http.StatusBadRequest,
			domain.ErrTitleTooShort,
		},
		{
			"too long title",
			args{
				ctx: ctx,
				req: &todo.CreateTodoRequest{
					Title: "a very long title that exceeds the maximum length of one hundred characters, which is not allowed in this test case............................................................................",
				},
			},
			http.StatusBadRequest,
			domain.ErrTitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, code, err := createTodoHandler.Handle(tt.args.ctx, tt.args.req)

			assert.Equal(t, code, tt.code)
			if err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

			} else {
				assert.NoError(t, err)
			}

		})
	}
}
