package user

import (
	"context"
	"net/http"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/sirupsen/logrus"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordResponse struct{}

type ForgotPasswordHandler struct {
	repo         Repository
	emailService MailService
	tokenService TokenService
	logger       *logrus.Logger
	validator    domain.Validator
}

func NewForgotPasswordHandler(repo Repository, emailService MailService, tokenService TokenService, logger *logrus.Logger, validator domain.Validator) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		repo:         repo,
		emailService: emailService,
		tokenService: tokenService,
		logger:       logger,
		validator:    validator,
	}
}

// Handle processes the request to initiate a password reset.
//
//	@Summary		Forgot Password
//	@Description	It sends  a password reset link to the user's email address.
//	@Tags			3- User
//	@Accept			json
//	@Produce		json
//	@Param			response-language	header	string					true	"Response Language"	enums(tr,ar,en)
//	@Param			request				body	ForgotPasswordRequest	true	"Forgot Password Request"
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/users/forgot-password [post]
func (h *ForgotPasswordHandler) Handle(ctx context.Context, req *ForgotPasswordRequest) (*ForgotPasswordResponse, int, error) {
	if err := h.validator.Validate(req); err != nil {
		return nil, http.StatusBadRequest, domain.ErrInvalidRequest
	}

	if exists, _ := h.repo.EmailExists(ctx, req.Email); !exists {
		return nil, http.StatusOK, nil
	}

	token, err := h.tokenService.GenerateTokenForForgotPassword(req.Email)
	if err != nil {
		h.logger.Error("failed to generate token for forgot password: ", err)
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	if err := h.emailService.SendPasswordResetEmail(
		req.Email,
		domain.ForgotPasswordEmailSubject,
		domain.NewForgotPasswordEmail(domain.NewForgotPasswordLink(token)),
	); err != nil {
		h.logger.Error("failed to send forgot password email: ", err)
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	return nil, http.StatusNoContent, nil
}
