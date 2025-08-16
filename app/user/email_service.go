package user

import (
	"context"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type MailService interface {
	SendSuccessfullyDeletedEmail(ctx context.Context, claims *domain.EmailClaims) error
	SendPasswordResetEmail(ctx context.Context, claims *domain.EmailClaims) error
	SendVerificationEmail(ctx context.Context, claims *domain.EmailClaims) error
}
