package testtodo

import (
	"context"
	"net/http"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	testUtils "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToggleCompletedTodoHandler(t *testing.T) {
	ctx := context.WithValue(context.Background(), domain.UserIDKey, domain.RealUserId)

	postgresContainer, connStr := testUtils.CreateTestContainer(t, ctx)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate postgres container")
	}()

	repo, err := postgresRepo.NewRepository(connStr)
	require.NoError(t, err, "failed to create repository")
	runMigrations(t, connStr)
	setupTestUser(t, connStr)
	setupTestTodo(t, connStr)

	updateTodoHandler := todo.NewToggleCompletedTodoHandler(repo)

	type args struct {
		ctx context.Context
		req *todo.ToggleCompletedTodoRequest
	}

	tests := []struct {
		name    string
		args    args
		code    int
		wantErr error
	}{
		{
			"valid toggle completed",
			args{
				ctx: ctx,
				req: &todo.ToggleCompletedTodoRequest{
					Id: domain.TestTodo.Id,
				},
			},
			http.StatusNoContent,
			nil,
		},
		{
			"invalid todo id",
			args{
				ctx: ctx,
				req: &todo.ToggleCompletedTodoRequest{
					Id: domain.FakeTodoUuid,
				},
			},
			http.StatusNotFound,
			domain.ErrTodoNotFound,
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
