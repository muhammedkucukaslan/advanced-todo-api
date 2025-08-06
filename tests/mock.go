package test

import (
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

// MockEmailService
type MockEmailService struct {
}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (m *MockEmailService) SendWelcomeEmail(name, to, subject, html string) error {
	return nil
}
func (m *MockEmailService) SendSuccessfullyDeletedEmail(to, email, subject, html string) error {
	return nil
}
func (m *MockEmailService) SendPasswordResetEmail(email, subject, html string) error {
	return nil
}
func (m *MockEmailService) SendVerificationEmail(name, to, subject, html string) error {
	return nil
}

// MokcValidator
type MockValidator struct{}

func NewMockValidator() *MockValidator {
	return &MockValidator{}
}

func (m *MockValidator) Validate(data any) error {
	return nil
}

// MockLogger
type MockLogger struct{}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (m *MockLogger) Info(msg string, args ...any)  {}
func (m *MockLogger) Error(msg string, args ...any) {}

// MockTokenService
type MockTokenService struct{}

func NewMockTokenService() *MockTokenService {
	return &MockTokenService{}
}

func (m *MockTokenService) GenerateToken(userID, role string, time time.Time) (string, error) {
	return domain.MockToken, nil
}

func (m *MockTokenService) ValidateToken(token string) (auth.TokenPayload, error) {
	return auth.TokenPayload{
		UserID: domain.TestUser.Id.String(),
		Role:   domain.TestUser.Role,
	}, nil
}

func (m *MockTokenService) GenerateVerificationToken(email string) (string, error) {
	return "mockedVerificationToken", nil
}

func (m *MockTokenService) ValidateVerifyEmailToken(tokenString string) (string, error) {
	return "", nil
}
func (m *MockTokenService) GenerateTokenForForgotPassword(email string) (string, error) {
	return "mockedForgotPasswordToken", nil
}

func (m *MockTokenService) ValidateForgotPasswordToken(tokenString string) (string, error) {
	return "", nil
}
