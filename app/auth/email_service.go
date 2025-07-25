package auth

type EmailService interface {
	SendWelcomeEmail(name, to, subject, html string) error
	SendVerificationEmail(from, to, subject, html string) error
}
