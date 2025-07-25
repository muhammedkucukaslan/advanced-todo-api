package user

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	"github.com/sirupsen/logrus"
)

type DeleteAccountRequest struct {
	Language string `reqHeader:"response-language" validate:"required,oneof=tr en ar"`
}

type DeleteAccountResponse struct{}

type DeleteAccountHandler struct {
	repo      Repository
	logger    *logrus.Logger
	validator *validator.Validate
	ms        MailService
}

func NewDeleteAccountHandler(repo Repository, logger *logrus.Logger, validate *validator.Validate, ms MailService) *DeleteAccountHandler {
	return &DeleteAccountHandler{repo: repo, logger: logger, ms: ms, validator: validate}
}

// Handle processes the request to delete a user's account.
// //	@Summary		Delete User Account
//
//	@Description	Delete a user's account
//	@Tags			3- User
//	@Param			response-language	header	string	true	"Response Language"	Enums(tr, en, ar)
//	@Security		BearerAuth
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/users/account [delete]
func (h *DeleteAccountHandler) Handle(ctx context.Context, req *DeleteAccountRequest) (*DeleteAccountResponse, int, error) {

	if err := h.validator.Struct(req); err != nil {
		return nil, 400, domain.ErrInvalidRequest
	}

	userId := domain.GetUserID(ctx)

	// deletes user current bucket and  user, also returns email
	fullName, email, err := h.repo.DeleteAccount(ctx, userId)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to delete user account: %v", err))
		return nil, 500, domain.ErrInternalServer
	}

	// TODO mail notification
	go func(fullName, email, language string) {
		maxRetries := 3
		retryInterval := 30 * time.Second

		for attempt := 1; attempt <= maxRetries; attempt++ {

			err := h.ms.SendSuccessfullyDeletedEmail(fullName, email,
				domain.SuccessfullyDeletedEmailSubject,
				domain.EnglishSuccessfullyDeletedEmail)

			if err == nil {
				return
			}
			h.logger.Errorf("attempt %d failed to send successfully deleted email to %s: %v", attempt, email, err)
			if attempt < maxRetries {
				time.Sleep(retryInterval)
			}
		}
		h.logger.Errorf("all %d attempts failed for successfully deleted email to %s", maxRetries, email)
		// TODO handle email sending failure (e.g., log it, notify admin, etc.)
	}(fullName, email, req.Language)

	return nil, 204, nil
}
