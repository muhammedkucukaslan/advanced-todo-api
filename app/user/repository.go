package user

import (
	"context"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*GetCurrentUserResponse, error)
	CheckEmail(ctx context.Context, email string) error
	UpdateAccount(ctx context.Context, id uuid.UUID, fullName, address, phone string) error
	DeleteAccount(ctx context.Context, id uuid.UUID) (string, string, error)
	GetUserOnlyHavingPasswordById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	EmailExists(ctx context.Context, email string) (bool, error)
	ResetPasswordByEmail(ctx context.Context, email, newPassword string) error
	ChangePassword(ctx context.Context, user *domain.User) error
	VerifyEmail(ctx context.Context, email string) error
	GetUserNameAndEmailByIdForSendingVerificationEmail(ctx context.Context, id uuid.UUID) (string, string, error)
	GetUserByIdForAdmin(ctx context.Context, id uuid.UUID) (*GetUserResponse, error)
	GetUsers(ctx context.Context, page, limit int) (GetUsersResponse, error)
}
