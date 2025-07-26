package testdomain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewTodo(t *testing.T) {
	type args struct {
		userId uuid.UUID
		title  string
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "should create valid todo",
			args: args{
				userId: uuid.New(),
				title:  "Buy groceries",
			},
			wantErr: nil,
		},
		{
			name: "should return error for empty title",
			args: args{
				userId: uuid.New(),
				title:  "",
			},
			wantErr: domain.ErrEmptyTitle,
		},
		{
			name: "should return error for title too long",
			args: args{
				userId: uuid.New(),
				title:  "a very long title that exceeds the maximum length of one hundred characters.........................................................................................",
			},
			wantErr: domain.ErrTitleTooLong,
		},
		{
			name: "should return error for title too short",
			args: args{
				userId: uuid.New(),
				title:  "ab",
			},
			wantErr: domain.ErrTitleTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewTodo(tt.args.userId, tt.args.title)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.NotEmpty(t, got.Id)
				assert.Equal(t, tt.args.userId, got.UserId)
				assert.Equal(t, tt.args.title, got.Title)
				assert.False(t, got.Completed)
				assert.WithinDuration(t, time.Now(), got.CreatedAt, time.Second)
				assert.Equal(t, time.Time{}, got.CompletedAt)
			}

		})
	}
}
