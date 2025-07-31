package testtodo

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodoHandler(t *testing.T) {

	repo, err := postgres.NewRepository(os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	createTodoHandler := todo.NewCreateTodoHandler(repo)
	ctx := context.WithValue(context.Background(), domain.UserIDKey, domain.RealUserId)
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
