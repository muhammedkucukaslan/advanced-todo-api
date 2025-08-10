package testtodo

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	mock "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodoHandler(t *testing.T) {
	ctx := context.WithValue(context.Background(), domain.UserIDKey, domain.RealUserId)

	createTodoHandler := todo.NewCreateTodoHandler(&MockRepository{}, mock.NewMockCache(), mock.NewMockLogger())

	validCreateTodoRequest := &todo.CreateTodoRequest{
		Title: "Test Todo",
	}

	invalidCreateTodoRequest := &todo.CreateTodoRequest{}

	tooShortCreateTodoRequest := &todo.CreateTodoRequest{
		Title: "ab",
	}

	tooLongCreateTodoRequest := &todo.CreateTodoRequest{
		Title: strings.Repeat("a", 105),
	}

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
		{"valid request", args{
			ctx: ctx,
			req: validCreateTodoRequest,
		}, http.StatusCreated, nil,
		},
		{"empty title", args{
			ctx: ctx,
			req: invalidCreateTodoRequest,
		}, http.StatusBadRequest, domain.ErrEmptyTitle,
		},
		{"too short title", args{
			ctx: ctx,
			req: tooShortCreateTodoRequest,
		}, http.StatusBadRequest, domain.ErrTitleTooShort,
		},
		{"too long title", args{
			ctx: ctx,
			req: tooLongCreateTodoRequest,
		}, http.StatusBadRequest, domain.ErrTitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, code, err := createTodoHandler.Handle(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.code, code)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
