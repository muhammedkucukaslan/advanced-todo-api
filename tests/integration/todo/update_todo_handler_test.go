package testtodo

import (
	"context"
	"net/http"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateTodoHandler(t *testing.T) {
	ctx := context.WithValue(context.Background(), domain.UserIDKey, domain.RealUserId)

	postgresContainer, connStr := createTestContainer(t, ctx)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate postgres container")
	}()

	repo, err := postgresRepo.NewRepository(connStr)
	require.NoError(t, err, "failed to create repository")
	runMigrations(t, connStr)
	setupTestUser(t, connStr)
	setupTestTodo(t, connStr)

	validator := validator.NewValidator(slog.NewLogger())
	updateTodoHandler := todo.NewUpdateTodoHandler(repo, validator)

	runMigrations(t, connStr)
	setupTestUser(t, connStr)

	type args struct {
		ctx context.Context
		req *todo.UpdateTodoRequest
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
				req: &todo.UpdateTodoRequest{
					Id:    domain.TestTodo.Id,
					Title: "Updated Test Todo",
				},
			},
			nil,
			http.StatusNoContent,
		},
		{
			"invalid request",
			args{
				ctx: ctx,
				req: &todo.UpdateTodoRequest{
					Id: domain.TestTodo.Id,
				},
			},
			domain.ErrInvalidRequest,
			http.StatusBadRequest,
		},
		{
			"too short title",
			args{
				ctx: ctx,
				req: &todo.UpdateTodoRequest{
					Id:    domain.TestTodo.Id,
					Title: "ab",
				},
			},
			domain.ErrTitleTooShort,
			http.StatusBadRequest,
		},
		{
			"too long title",
			args{
				ctx: ctx,
				req: &todo.UpdateTodoRequest{
					Id:    domain.TestTodo.Id,
					Title: "a very long title that exceeds the maximum length of one hundred characters, which is not allowed in this test case............................................................................",
				},
			},
			domain.ErrTitleTooLong,
			http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, code, err := updateTodoHandler.Handle(tt.args.ctx, tt.args.req)

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
