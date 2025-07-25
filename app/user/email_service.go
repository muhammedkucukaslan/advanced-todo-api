package user

type MailService interface {
	SendSuccessfullyDeletedEmail(to, email, subject, html string) error
	SendPasswordResetEmail(email, subject, html string) error
	SendVerificationEmail(to, email, subject, html string) error
}
