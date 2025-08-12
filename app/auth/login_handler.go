package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginHandler struct {
	repo      Repository
	ts        TokenService
	validator domain.Validator
	logger    domain.Logger
}

func NewLoginHandler(repo Repository, ts TokenService, validator domain.Validator, logger domain.Logger) *LoginHandler {
	return &LoginHandler{repo: repo, ts: ts, validator: validator, logger: logger}
}

// @Summary		Login
// @Description	Login a user or admin
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		LoginRequest	true	"Login Request"
// @Success		200		{object}	LoginResponse
// @Failure		400
// @Failure		404
// @Failure		500
// @Router			/login [post]
func (h *LoginHandler) Handle(ctx context.Context, req *LoginRequest) (*LoginResponse, int, error) {
	if err := h.validator.Validate(req); err != nil {
		return nil, http.StatusBadRequest, domain.ErrInvalidRequest
	}

	user, err := h.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {

		if errors.Is(err, domain.ErrEmailNotFound) {
			return nil, http.StatusNotFound, domain.ErrEmailNotFound
		}
		h.logger.Error("error while getting user by email: ", err)
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	if err = user.ValidatePassword(req.Password); err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return nil, http.StatusBadRequest, domain.ErrInvalidCredentials
		}
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	token, err := h.ts.GenerateAuthToken(user.Id.String(), user.Role)
	if err != nil {
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	return &LoginResponse{Token: token}, http.StatusOK, nil
}
