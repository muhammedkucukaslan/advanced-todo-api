package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type SignupRequest struct {
	FullName string `json:"fullName" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type SignupResponse struct {
	AccessToken string `json:"access_token"`
}

type SignupConfig struct {
	RefreshTokenCookieDuration time.Duration
	Secure                     bool
	Repo                       Repository
	TokenService               TokenService
	CookieService              CookieService
	EmailService               EmailService
	Validator                  domain.Validator
	Logger                     domain.Logger
}

type SignupHandler struct {
	refreshTokenCookieDuration time.Duration
	secure                     bool
	repo                       Repository
	cs                         CookieService
	ts                         TokenService
	es                         EmailService
	validator                  domain.Validator
	logger                     domain.Logger
}

func NewSignupHandler(config *SignupConfig) *SignupHandler {
	return &SignupHandler{
		refreshTokenCookieDuration: config.RefreshTokenCookieDuration,
		secure:                     config.Secure,
		repo:                       config.Repo,
		cs:                         config.CookieService,
		ts:                         config.TokenService,
		es:                         config.EmailService,
		validator:                  config.Validator,
		logger:                     config.Logger,
	}
}

// @Summary		Signup
// @Description	Signup a new user
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		SignupRequest	true	"Signup request"
// @Success		200		{object}	SignupResponse
// @Failure		400
// @Failure		409
// @Failure		500
// @Router			/signup [post]
func (h *SignupHandler) Handle(ctx context.Context, req *SignupRequest) (*SignupResponse, int, error) {
	if err := h.validator.Validate(req); err != nil {
		return nil, http.StatusBadRequest, domain.ErrInvalidRequest
	}

	user, err := domain.NewUser(req.FullName, req.Password, req.Email)
	if err != nil {
		if !errors.Is(err, domain.ErrInternalServer) {
			return nil, http.StatusBadRequest, err
		}
		h.logger.Error("error while creating domain user: ", err)
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	accessToken, err := h.ts.GenerateAuthAccessToken(user.Id.String(), user.Role)
	if err != nil {
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	refreshToken, err := h.ts.GenerateAuthRefreshToken(user.Id.String(), user.Role)
	if err != nil {
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	err = h.repo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return nil, http.StatusConflict, domain.ErrEmailAlreadyExists
		}
		h.logger.Error("error while creating user in repository: ", err)
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	if err := h.repo.SaveRefreshToken(ctx, domain.NewRefreshToken(
		user.Id,
		refreshToken,
		h.refreshTokenCookieDuration,
	)); err != nil {
		h.logger.Error("error while saving refresh token: ", err)
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	go func(fullname, email string) {
		var err error
		const maxRetries = 3
		const retryInterval = 30 * time.Second

		for attempt := 1; attempt <= maxRetries; attempt++ {
			verificationToken, err := h.ts.GenerateEmailVerificationToken(email)
			if err != nil {
				h.logger.Error(fmt.Sprintf("[Signup] Attempt %d: Failed to generate token for %s: %v", attempt, email, err))
				time.Sleep(retryInterval)
				continue
			}

			err = h.es.SendVerificationEmail(
				fullname,
				email,
				domain.VerificationEmailSubject,
				domain.NewVerificationEmailBody(domain.NewVerificationEmailLink(verificationToken)),
			)

			if err == nil {
				h.logger.Info(fmt.Sprintf("[Signup] Verification email sent to %s", email))
				return
			}

			h.logger.Error(fmt.Sprintf("[Signup] Attempt %d: Failed to send email to %s: %v", attempt, email, err))
			time.Sleep(retryInterval)
		}

		h.logger.Error(fmt.Sprintf("[Signup] All retries failed for %s: %v", email, err))
		// TODO: notify admin, push to dead-letter queue, etc.
	}(user.FullName, user.Email)

	go func(fullname, email string) {
		maxRetries := 3
		retryInterval := 30 * time.Second

		for attempt := 1; attempt <= maxRetries; attempt++ {

			err := h.es.SendWelcomeEmail(fullname, email,
				domain.WelcomeEmailSubject,
				domain.NewWelcomeEmailBody(fullname))

			if err == nil {
				return
			}

			if attempt < maxRetries {
				time.Sleep(retryInterval)
			}
		}
		h.logger.Error(fmt.Sprintf("all %d attempts failed for welcome email to %s", maxRetries, email))
		// TODO handle email sending failure (e.g., log it, notify admin, etc.)
	}(user.FullName, user.Email)

	h.cs.SetRefreshToken(ctx, &RefreshTokenCookieClaims{
		Token:    refreshToken,
		Duration: h.refreshTokenCookieDuration,
		Secure:   h.secure,
	})

	return &SignupResponse{AccessToken: accessToken}, http.StatusCreated, nil
}
