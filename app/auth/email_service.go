package auth

import "context"

type MailClaims struct {
	Name    string
	To      string
	Subject string
	HTML    string
}

type EmailService interface {
	SendWelcomeEmail(ctx context.Context, ml *MailClaims) error
	SendVerificationEmail(ctx context.Context, ml *MailClaims) error
}
