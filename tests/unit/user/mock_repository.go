package testuser

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/user"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type MockRepository struct{}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) GetUserById(ctx context.Context, id uuid.UUID) (*user.GetCurrentUserResponse, error) {
	// mock değer döndür
	return &user.GetCurrentUserResponse{
		Id:       id.String(),
		FullName: "Mock User",
		Email:    "mock@example.com",
	}, nil
}

func (m *MockRepository) CheckEmail(ctx context.Context, email string) error {
	// email geçerliyse nil döndür, değilse error
	if email == "exists@example.com" {
		return errors.New("email already exists")
	}
	return nil
}

func (m *MockRepository) DeleteAccount(ctx context.Context, id uuid.UUID) (string, string, error) {
	return "Mock User", "mock@example.com", nil
}

func (m *MockRepository) GetUserOnlyHavingPasswordById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return &domain.User{
		Id:       id,
		Password: "$2y$10$7ALFQvvizAtvcM.zmnoZHOBDAPVQfrxJ4gPf/vzFzio.zPWYlFE5W", // hash for "password123"
	}, nil
}

func (m *MockRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	if email == "exists@example.com" {
		return true, nil
	}
	return false, nil
}

func (m *MockRepository) ResetPasswordByEmail(ctx context.Context, email, newPassword string) error {
	return nil
}

func (m *MockRepository) ChangePassword(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *MockRepository) UpdateFullName(ctx context.Context, id uuid.UUID, fullName string) error {
	return nil
}

func (m *MockRepository) VerifyEmail(ctx context.Context, email string) error {
	return nil
}

func (m *MockRepository) GetUserNameAndEmailByIdForSendingVerificationEmail(ctx context.Context, id uuid.UUID) (string, string, error) {
	return "Mock User", "mock@example.com", nil
}

func (m *MockRepository) GetUserByIdForAdmin(ctx context.Context, id uuid.UUID) (*user.GetUserResponse, error) {
	return &user.GetUserResponse{
		ID:       id,
		FullName: "Admin View User",
		Email:    "adminview@example.com",
	}, nil
}

func (m *MockRepository) GetUsers(ctx context.Context, page, limit int) (user.GetUsersResponse, error) {
	users := user.GetUsersResponse{
		{
			Id:       uuid.New(),
			FullName: "User One",
			Email:    "user1@example.com",
		},
		{
			Id:       uuid.New(),
			FullName: "User Two",
			Email:    "user2@example.com",
		},
	}

	return users, nil
}
