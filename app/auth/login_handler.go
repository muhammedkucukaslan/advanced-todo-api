package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type LoginConfig struct {
	RefreshTokenCookieDuration time.Duration
	Secure                     bool
	Repo                       Repository
	TokenService               TokenService
	CookieService              CookieService
	Validator                  domain.Validator
	Logger                     domain.Logger
}

type LoginHandler struct {
	refreshTokenCookieDuration time.Duration
	secure                     bool
	repo                       Repository
	ts                         TokenService
	validator                  domain.Validator
	logger                     domain.Logger
	cs                         CookieService
}

func NewLoginHandler(config *LoginConfig) *LoginHandler {
	return &LoginHandler{
		refreshTokenCookieDuration: config.RefreshTokenCookieDuration,
		secure:                     config.Secure,
		repo:                       config.Repo,
		ts:                         config.TokenService,
		validator:                  config.Validator,
		logger:                     config.Logger,
		cs:                         config.CookieService,
	}
}

// @Summary		Login
// @Description	Login a user or admin
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		LoginRequest	true	"Login Request"
// @Success		200		{object}	LoginResponse
// @Failure		400
// @Failure		401
// @Failure		404
// @Failure		500
// @Router			/auth/login [post]
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

	token, err := h.ts.GenerateAuthAccessToken(user.Id.String(), user.Role)
	if err != nil {
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	// TODO delete old refresh token
	refreshToken, err := h.ts.GenerateAuthRefreshToken(user.Id.String(), user.Role)
	if err != nil {
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	if err := h.repo.UpsertRefreshToken(ctx, domain.NewRefreshToken(
		user.Id,
		refreshToken,
		h.refreshTokenCookieDuration,
	)); err != nil {
		h.logger.Error("error while saving refresh token: ", err)
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	h.cs.SetRefreshToken(ctx, &RefreshTokenCookieClaims{
		Token:    refreshToken,
		Duration: h.refreshTokenCookieDuration,
		Secure:   h.secure,
	})

	return &LoginResponse{AccessToken: token}, http.StatusOK, nil
}
