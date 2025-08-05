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

func TestGetTodoHandler(t *testing.T) {
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

	getTodoByIdHandler := todo.NewGetTodoByIdHandler(repo)

	type args struct {
		ctx context.Context
		req *todo.GetTodoByIdRequest
	}

	tests := []struct {
		name    string
		args    args
		want    *todo.GetTodoByIdResponse
		code    int
		wantErr error
	}{
		{"valid request", args{
			ctx: ctx,
			req: &todo.GetTodoByIdRequest{
				Id: domain.TestTodo.Id,
			},
		}, &todo.GetTodoByIdResponse{
			Id:        domain.TestTodo.Id,
			Title:     domain.TestTodo.Title,
			Completed: domain.TestTodo.Completed,
		}, http.StatusOK, nil},
		{"invalid todo id", args{
			ctx: ctx,
			req: &todo.GetTodoByIdRequest{
				Id: domain.FakeTodoIdUuid,
			},
		}, nil, http.StatusNotFound, domain.ErrTodoNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, code, err := getTodoByIdHandler.Handle(tt.args.ctx, tt.args.req)

			assert.Equal(t, code, tt.code)
			if err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.Title, got.Title)
				assert.Equal(t, tt.want.Completed, got.Completed)
			}
		})
	}
}
