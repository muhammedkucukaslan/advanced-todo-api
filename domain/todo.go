package domain

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	UserId      uuid.UUID
	Id          uuid.UUID
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

func NewTodo(userId uuid.UUID, title string) (*Todo, error) {
	if len(title) == 0 {
		return nil, ErrEmptyTitle
	}

	if userId == uuid.Nil {
		return nil, ErrUserIdCannotBeEmpty
	}

	if len(title) > 100 {
		return nil, ErrTitleTooLong
	}

	if len(title) < 3 {
		return nil, ErrTitleTooShort
	}

	return &Todo{
		UserId:      userId,
		Id:          uuid.New(),
		Title:       title,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}, nil
}
