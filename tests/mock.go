package test

import (
	"context"
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

func (m *MockTokenService) GenerateAuthAccessToken(userID, role string) (string, error) {
	return domain.MockToken, nil
}

func (m *MockTokenService) ValidateAuthAccessToken(token string) (*auth.TokenPayload, error) {
	return &auth.TokenPayload{
		UserID: domain.TestUser.Id.String(),
		Role:   domain.TestUser.Role,
	}, nil
}

func (m *MockTokenService) GenerateEmailVerificationToken(email string) (string, error) {
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

// MockCache
type MockCache struct {
}

func NewMockCache() *MockCache {
	return &MockCache{}
}

func (m *MockCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (m *MockCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return nil
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	return nil
}
