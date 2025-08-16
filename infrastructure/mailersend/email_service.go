package mailsend

import (
	"context"

	"github.com/mailersend/mailersend-go"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type MailerSendService struct {
	client      *mailersend.Mailersend
	SenderEmail string
	SenderName  string
}

func NewMailerSendService(apiKey, senderEmail, senderName string) *MailerSendService {
	return &MailerSendService{client: mailersend.NewMailersend(apiKey), SenderEmail: senderEmail, SenderName: senderName}
}

func (m *MailerSendService) SendWelcomeEmail(ctx context.Context, claims *domain.EmailClaims) error {

	from := mailersend.From{
		Name:  m.SenderName,
		Email: m.SenderEmail,
	}

	recipients := []mailersend.Recipient{
		{
			Name:  claims.Name,
			Email: claims.To,
		},
	}

	message := m.client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(claims.Subject)
	message.SetHTML(claims.HTML)

	_, err := m.client.Email.Send(context.Background(), message)
	return err
}

func (m *MailerSendService) SendSuccessfullyDeletedEmail(ctx context.Context, claims *domain.EmailClaims) error {

	from := mailersend.From{
		Name:  m.SenderName,
		Email: m.SenderEmail,
	}

	recipients := []mailersend.Recipient{
		{
			Name:  claims.Name,
			Email: claims.To,
		},
	}

	message := m.client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(claims.Subject)
	message.SetHTML(claims.HTML)

	_, err := m.client.Email.Send(context.Background(), message)
	return err
}

func (m *MailerSendService) SendPasswordResetEmail(email, subject, html string) error {

	from := mailersend.From{
		Name:  m.SenderName,
		Email: m.SenderEmail,
	}

	recipients := []mailersend.Recipient{
		{
			Email: email,
		},
	}

	message := m.client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)

	_, err := m.client.Email.Send(context.Background(), message)
	return err
}

func (m *MailerSendService) SendVerificationEmail(ctx context.Context, claims *domain.EmailClaims) error {

	from := mailersend.From{
		Name:  m.SenderName,
		Email: m.SenderEmail,
	}

	recipients := []mailersend.Recipient{
		{
			Name:  claims.Name,
			Email: claims.To,
		},
	}

	message := m.client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(claims.Subject)
	message.SetHTML(claims.HTML)

	_, err := m.client.Email.Send(ctx, message)
	return err
}
