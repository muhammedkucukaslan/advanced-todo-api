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

func (m *MockEmailService) SendWelcomeEmail(ctx context.Context, claims *domain.EmailClaims) error {
	return nil
}
func (m *MockEmailService) SendSuccessfullyDeletedEmail(ctx context.Context, claims *domain.EmailClaims) error {
	return nil
}
func (m *MockEmailService) SendPasswordResetEmail(ctx context.Context, claims *domain.EmailClaims) error {
	return nil
}
func (m *MockEmailService) SendVerificationEmail(ctx context.Context, claims *domain.EmailClaims) error {
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

func (m *MockTokenService) GenerateAuthRefreshToken(userID, role string) (string, error) {
	return "mockedRefreshToken", nil
}

func (m *MockTokenService) ValidateAuthRefreshToken(token string) (*auth.TokenPayload, error) {
	return &auth.TokenPayload{
		UserID: domain.TestUser.Id.String(),
		Role:   domain.TestUser.Role,
	}, nil
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

// MockCookie Service

type MockCookieService struct {
}

func NewMockCookieService() *MockCookieService {
	return &MockCookieService{}
}

func (m *MockCookieService) SetRefreshToken(ctx context.Context, token string) {
}

func (m *MockCookieService) SetAccessToken(ctx context.Context, token string) {
}

func (m *MockCookieService) RemoveTokens(ctx context.Context) {

}
