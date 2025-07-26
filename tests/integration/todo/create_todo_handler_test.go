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

	// This is a real user ID for testing
	realUserId := "9796b492-662d-4742-8f7b-5eefbe1da107"
	mockUserId := "5ee1903d-0c9a-4d95-aae2-7215e168564b"

	createTodoHandler := todo.NewCreateTodoHandler(repo)
	ctx := context.WithValue(context.Background(), domain.UserIDKey, realUserId)
	ctxWithFakeUserId := context.WithValue(context.Background(), domain.UserIDKey, mockUserId)

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
			"should create todo successfully",
			args{
				ctx: ctx,
				req: validCreateTodoRequest,
			},
			nil,
			http.StatusCreated,
		},
		{
			"should return forbidden error for userId which is not in database",
			args{
				ctx: ctxWithFakeUserId,
				req: validCreateTodoRequest,
			},
			domain.ErrUserNotFound,
			http.StatusForbidden,
		},

		{
			"should return bad request error for invalid title",
			args{
				ctx: ctx,
				req: invalidCreateTodoRequest,
			},
			domain.ErrInvalidRequest,
			http.StatusBadRequest,
		},
		{
			"should return bad request error for too short title",
			args{
				ctx: ctx,
				req: tooShortCreateTodoRequest,
			},
			domain.ErrTitleTooShort,
			http.StatusBadRequest,
		},
		{
			"should return bad request error for too long title",
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

			if err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, code, tt.code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, code, tt.code)
			}

		})
	}
}
