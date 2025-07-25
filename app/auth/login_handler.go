package auth

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/sirupsen/logrus"
)

type LoginRequest struct {
	Language string `reqHeader:"response-language" validate:"required,oneof=en tr ar" swaggerignore:"true"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginHandler struct {
	repo      Repository
	ts        TokenService
	validator *validator.Validate
	logger    *logrus.Logger
}

func NewLoginHandler(repo Repository, ts TokenService, validator *validator.Validate, logger *logrus.Logger) *LoginHandler {
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
	if err := h.validator.Struct(req); err != nil {
		return nil, 400, domain.ErrInvalidRequest
	}

	user, err := h.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {

		if errors.Is(err, domain.ErrEmailNotFound) {
			return nil, 404, domain.ErrEmailNotFound
		}
		h.logger.Error("error while getting user by email: ", err)
		return nil, 500, domain.ErrInternalServer
	}

	err = user.ValidatePassword(req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return nil, 400, domain.ErrInvalidCredentials
		}
		h.logger.Error("error while validating password: ", err)
		return nil, 500, domain.ErrInternalServer
	}

	token, err := h.ts.GenerateToken(user.Id.String(), user.Role, time.Now())
	if err != nil {
		h.logger.Error("error while generating token: ", err)
		return nil, 500, domain.ErrInternalServer
	}

	return &LoginResponse{Token: token}, 200, nil
}
