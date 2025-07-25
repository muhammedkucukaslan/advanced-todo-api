package user

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type SendVerificationEmailRequest struct {
	Language string `reqHeader:"response-language" validate:"required,oneof=tr en ar" swaggerignore:"true"`
}

type SendVerificationEmailResponse struct{}

type SendVerificationEmailHandler struct {
	repo         Repository
	validate     *validator.Validate
	tokenService TokenService
	ms           MailService
}

func NewSendVerificationEmailHandler(repo Repository, validate *validator.Validate, tokenService TokenService, ms MailService) *SendVerificationEmailHandler {
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
//	@Tags			3- User
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			response-language	header	string	true	"Response Language"	Enums(tr, en, ar)
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/users/send-verification-email [post]
func (h *SendVerificationEmailHandler) Handle(ctx context.Context, req *SendVerificationEmailRequest) (*SendVerificationEmailResponse, int, error) {

	if err := h.validate.Struct(req); err != nil {
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

	token, err := h.tokenService.GenerateVerificationToken(email)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err = h.ms.SendVerificationEmail(
		fullname,
		email,
		domain.VerificationEmailSubject,
		domain.NewVerificationEmailBody(domain.NewVerificationEmailLink(token)),
	); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusNoContent, nil
}
