package user

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/sirupsen/logrus"
)

type ForgotPasswordRequest struct {
	Language string `reqHeader:"response-language" validate:"required,oneof=tr en ar" swaggerignore:"true"`
	Email    string `json:"email" validate:"required,email"`
}

type ForgotPasswordResponse struct{}

type ForgotPasswordHandler struct {
	repo         Repository
	emailService MailService
	tokenService TokenService
	logger       *logrus.Logger
	validator    *validator.Validate
}

func NewForgotPasswordHandler(repo Repository, emailService MailService, tokenService TokenService, logger *logrus.Logger, validator *validator.Validate) *ForgotPasswordHandler {
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
	if err := h.validator.Struct(req); err != nil {
		return nil, 400, domain.ErrInvalidRequest
	}

	if exists, _ := h.repo.EmailExists(ctx, req.Email); !exists {
		return nil, 200, nil
	}

	token, err := h.tokenService.GenerateTokenForForgotPassword(req.Email)
	if err != nil {
		h.logger.Error("failed to generate token for forgot password: ", err)
		return nil, 500, domain.ErrInternalServer
	}

	if err := h.emailService.SendPasswordResetEmail(
		req.Email,
		domain.NewForgotPasswordSubject(req.Language),
		domain.NewForgotPasswordEmail(domain.NewForgotPasswordLink(token), req.Language),
	); err != nil {
		h.logger.Error("failed to send forgot password email: ", err)
		return nil, 500, domain.ErrInternalServer
	}

	return nil, 204, nil
}
