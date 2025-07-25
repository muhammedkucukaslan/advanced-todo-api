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

func NewTodo(userId uuid.UUID, title string) *Todo {
	return &Todo{
		UserId:    userId,
		Id:        uuid.New(),
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
}
