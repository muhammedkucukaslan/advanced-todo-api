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
	Token string `json:"token"`
}

type SignupHandler struct {
	repo      Repository
	ts        TokenService
	es        EmailService
	validator domain.Validator
	logger    domain.Logger
}

func NewSignupHandler(repo Repository, ts TokenService, es EmailService, validator domain.Validator, logger domain.Logger) *SignupHandler {
	return &SignupHandler{repo: repo, ts: ts, es: es, validator: validator, logger: logger}
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

	token, err := h.ts.GenerateAuthToken(user.Id.String(), user.Role)
	if err != nil {
		h.logger.Error("error while generating token: ", err)
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

	go func(fullname, email string) {
		var err error
		const maxRetries = 3
		const retryInterval = 30 * time.Second

		for attempt := 1; attempt <= maxRetries; attempt++ {
			verificationToken, err := h.ts.GenerateVerificationToken(email)
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

	return &SignupResponse{Token: token}, http.StatusCreated, nil
}
