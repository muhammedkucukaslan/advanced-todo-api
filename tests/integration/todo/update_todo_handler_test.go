package integrationtest_todo

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	testUtils "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateTodoHandler(t *testing.T) {
	t.Parallel()
	ctx := context.WithValue(context.Background(), domain.UserIDKey, domain.RealUserId)

	postgresContainer, connStr := testUtils.CreatePostgresTestContainer(t, ctx)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate postgres container")
	}()

	repo := postgresRepo.NewRepository(connStr)
	runMigrations(t, connStr)
	setupTestUser(t, connStr)
	setupTestTodo(t, connStr)

	updateTodoHandler := todo.NewUpdateTodoHandler(repo)

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
			"valid update",
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
		{"not found", args{
			ctx: ctx,
			req: &todo.UpdateTodoRequest{
				Id:    domain.FakeTodoUuid,
				Title: "Updated Test Todo",
			},
		}, domain.ErrTodoNotFound, http.StatusNotFound},
		{
			"empty title",
			args{
				ctx: ctx,
				req: &todo.UpdateTodoRequest{
					Id: domain.TestTodo.Id,
				},
			},
			domain.ErrEmptyTitle,
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
					Title: strings.Repeat("a", 101),
				},
			},
			domain.ErrTitleTooLong,
			http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, code, err := updateTodoHandler.Handle(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.code, code)
			if err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
