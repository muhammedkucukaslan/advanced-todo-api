package auth

import (
	"context"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type EmailService interface {
	SendWelcomeEmail(ctx context.Context, ml *domain.EmailClaims) error
	SendVerificationEmail(ctx context.Context, ml *domain.EmailClaims) error
}
