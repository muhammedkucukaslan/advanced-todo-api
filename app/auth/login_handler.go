package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

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
// @Tags			2- Auth
// @Accept			json
// @Produce		json
// @Param			response-language	header		string			true	"Response Language"	Enums(tr, ar, en)
// @Param			request				body		LoginRequest	true	"Login Request"
// @Success		200					{object}	LoginResponse
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
	err = user.ValidatePassword(req.Password)
	if err != nil {
		fmt.Println("Password validation error:", err)
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return nil, http.StatusBadRequest, domain.ErrInvalidCredentials
		}
		h.logger.Error("error while validating password: ", err)
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	token, err := h.ts.GenerateToken(user.Id.String(), user.Role, time.Now())
	if err != nil {
		h.logger.Error("error while generating token: ", err)
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	return &LoginResponse{Token: token}, http.StatusOK, nil
}
