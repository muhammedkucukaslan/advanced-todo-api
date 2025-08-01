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

	runMigrations(t, connStr)
	setupTestUser(t, connStr)

	ctxWithFakeUserId := context.WithValue(context.Background(), domain.UserIDKey, domain.FakeUserId)

	validCreateTodoRequest := &todo.CreateTodoRequest{
		Title: "Test Todo",
	}

	invalidCreateTodoRequest := &todo.CreateTodoRequest{}

	tooShortCreateTodoRequest := &todo.CreateTodoRequest{
		Title: "ab",
	}

	tooLongCreateTodoRequest := &todo.CreateTodoRequest{
		Title: "a very long title that exceeds the maximum length of one hundred characters, which is not allowed in this test case............................................................................",
	}

	type args struct {
		ctx context.Context
		req *todo.CreateTodoRequest
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
		code    int
	}{
		{
			"valid creation",
			args{
				ctx: ctx,
				req: validCreateTodoRequest,
			},
			nil,
			http.StatusCreated,
		},
		{
			"invalid user ID request",
			args{
				ctx: ctxWithFakeUserId,
				req: validCreateTodoRequest,
			},
			domain.ErrUserNotFound,
			http.StatusForbidden,
		},

		{
			"invalid request",
			args{
				ctx: ctx,
				req: invalidCreateTodoRequest,
			},
			domain.ErrInvalidRequest,
			http.StatusBadRequest,
		},
		{
			"too short title",
			args{
				ctx: ctx,
				req: tooShortCreateTodoRequest,
			},
			domain.ErrTitleTooShort,
			http.StatusBadRequest,
		},
		{
			"too long title",
			args{
				ctx: ctx,
				req: tooLongCreateTodoRequest,
			},
			domain.ErrTitleTooLong,
			http.StatusBadRequest,
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
