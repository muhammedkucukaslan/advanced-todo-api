package testtodo

import (
	"context"
	"net/http"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodoHandlerWithMock(t *testing.T) {
	ctx := context.WithValue(context.Background(), domain.UserIDKey, "5ee1903d-0c9a-4d95-aae2-7215e168564b")

	createTodoHandler := todo.NewCreateTodoHandler(&MockRepository{})

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
			"invalid title creation",
			args{
				ctx: ctx,
				req: invalidCreateTodoRequest,
			},
			domain.ErrInvalidRequest,
			http.StatusBadRequest,
		},
		{
			"too short title creation",
			args{
				ctx: ctx,
				req: tooShortCreateTodoRequest,
			},
			domain.ErrTitleTooShort,
			http.StatusBadRequest,
		},
		{
			"too long title creation",
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
