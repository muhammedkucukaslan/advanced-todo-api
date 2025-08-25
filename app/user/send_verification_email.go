package user

import (
	"context"
	"net/http"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type SendVerificationEmailRequest struct {
}

type SendVerificationEmailResponse struct{}

type SendVerificationEmailHandler struct {
	repo         Repository
	validate     domain.Validator
	tokenService TokenService
	ms           MailService
}

func NewSendVerificationEmailHandler(repo Repository, validate domain.Validator, tokenService TokenService, ms MailService) *SendVerificationEmailHandler {
	return &SendVerificationEmailHandler{
		repo:         repo,
		validate:     validate,
		tokenService: tokenService,
		ms:           ms,
	}
}

// Handle processes the request to send a verification email to the user.
//
//	@Summary		Send Verification Email
//
//	@Description	Sends a verification email to the user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/users/send-verification-email [post]
func (h *SendVerificationEmailHandler) Handle(ctx context.Context, req *SendVerificationEmailRequest) (*SendVerificationEmailResponse, int, error) {

	if err := h.validate.Validate(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	userId := domain.GetUserID(ctx)

	fullname, email, err := h.repo.GetUserNameAndEmailByIdForSendingVerificationEmail(ctx, userId)
	if err != nil {
		if err == domain.ErrEmailAlreadyVerified {
			return nil, http.StatusBadRequest, domain.ErrEmailAlreadyVerified
		}
		return nil, http.StatusInternalServerError, err
	}

	token, err := h.tokenService.GenerateSecureEmailToken(email)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err = h.ms.SendVerificationEmail(context.Background(), &domain.EmailClaims{
		Name:    fullname,
		To:      email,
		Subject: domain.VerificationEmailSubject,
		HTML:    domain.NewVerificationEmailBody(domain.NewVerificationEmailLink(token)),
	}); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusNoContent, nil
}
