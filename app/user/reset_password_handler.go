package user

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/sirupsen/logrus"
)

type ResetPasswordRequest struct {
	Language string `reqHeader:"response-language" validate:"required,oneof=tr en ar" swaggerignore:"true"`
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type ResetPasswordResponse struct{}

type ResetPasswordHandler struct {
	repo         Repository
	tokenService TokenService
	logger       *logrus.Logger
	validator    *validator.Validate
}

func NewResetPasswordHandler(repo Repository, tokenService TokenService, logger *logrus.Logger, validator *validator.Validate) *ResetPasswordHandler {
	return &ResetPasswordHandler{
		repo:         repo,
		tokenService: tokenService,
		logger:       logger,
		validator:    validator,
	}
}

// Handle processes the request to reset a user's password using a token.
//
//	@Summary		Reset Password
//	@Description	It resets a user's password using a token.
//	@Tags			3- User
//	@Accept			json
//	@Produce		json
//	@Param			response-language	header	string					true	"Response Language"	Enums(tr, en, ar)
//	@Param			request				body	ResetPasswordRequest	true	"Reset Password Request"
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/users/reset-password [post]
func (h *ResetPasswordHandler) Handle(ctx context.Context, req *ResetPasswordRequest) (*ResetPasswordResponse, int, error) {
	if err := h.validator.Struct(req); err != nil {
		return nil, 400, domain.ErrInvalidRequest
	}

	email, err := h.tokenService.ValidateForgotPasswordToken(req.Token)
	if err != nil {
		h.logger.Error("failed to validate token for forgot password: ", err)
		return nil, 401, domain.ErrUnauthorized
	}

	hashedPassword, err := domain.HashPassword(req.Password)
	if err != nil {
		return nil, 500, domain.ErrInternalServer
	}

	if err := h.repo.ResetPasswordByEmail(ctx, email, hashedPassword); err != nil {
		return nil, 500, domain.ErrInternalServer
	}
	return nil, 204, nil
}
