package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type ResetPasswordResponse struct{}

type ResetPasswordHandler struct {
	repo         Repository
	tokenService TokenService
	logger       domain.Logger
	validator    domain.Validator
}

func NewResetPasswordHandler(repo Repository, tokenService TokenService, logger domain.Logger, validator domain.Validator) *ResetPasswordHandler {
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
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request	body	ResetPasswordRequest	true	"Reset Password Request"
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/users/reset-password [post]
func (h *ResetPasswordHandler) Handle(ctx context.Context, req *ResetPasswordRequest) (*ResetPasswordResponse, int, error) {
	if err := h.validator.Validate(req); err != nil {
		return nil, http.StatusBadRequest, domain.ErrInvalidRequest
	}

	email, err := h.tokenService.ValidateForgotPasswordToken(req.Token)
	if err != nil {
		h.logger.Error("failed to validate token for forgot password: ", err)
		return nil, http.StatusUnauthorized, domain.ErrUnauthorized
	}

	hashedPassword, err := domain.HashPassword(req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrPasswordTooShort) {
			return nil, http.StatusBadRequest, domain.ErrPasswordTooShort
		}
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	if err := h.repo.ResetPasswordByEmail(ctx, email, hashedPassword); err != nil {

		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}
	return nil, http.StatusNoContent, nil
}
