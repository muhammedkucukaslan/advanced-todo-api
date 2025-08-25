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

func (s *MockTokenService) GenerateAuthAccessToken(userID string, role string) (string, error) {
	return "mockAccessToken", nil
}

func (s *MockTokenService) GenerateAuthRefreshToken(userID string, role string) (string, error) {
	return "mockRefreshToken", nil
}

func (s *MockTokenService) GenerateSecureEmailToken(email string) (string, error) {
	return "mockEmailToken", nil
}

func (s *MockTokenService) ValidateAuthAccessToken(tokenString string) (*auth.TokenPayload, error) {
	return &auth.TokenPayload{
		UserID: "mockUserID",
		Role:   "mockRole",
	}, nil
}

func (s *MockTokenService) ValidateAuthRefreshToken(tokenString string) (*auth.TokenPayload, error) {
	return &auth.TokenPayload{
		UserID: "mockUserID",
		Role:   "mockRole",
	}, nil
}

func (s *MockTokenService) ValidateSecureEmailToken(tokenString string) (string, error) {
	return "mockEmail", nil
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
