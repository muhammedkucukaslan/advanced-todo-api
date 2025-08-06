package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	MaxTitleLength = 100
	MinTitleLength = 3
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

	if IsUserIdEmpty(userId) {
		return nil, ErrUserIdCannotBeEmpty
	}

	if err := ValidateTitle(title); err != nil {
		return nil, err
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

func IsUserIdEmpty(userId uuid.UUID) bool {
	return userId == uuid.Nil
}

func ValidateTitle(title string) error {
	if len(title) == 0 {
		return ErrEmptyTitle
	}
	if len(title) > MaxTitleLength {
		return ErrTitleTooLong
	}
	if len(title) < MinTitleLength {
		return ErrTitleTooShort
	}
	return nil
}
