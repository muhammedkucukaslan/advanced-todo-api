package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type DeleteAccountRequest struct {
}

type DeleteAccountResponse struct{}

type DeleteAccountHandler struct {
	repo   Repository
	logger domain.Logger
	ms     MailService
}

func NewDeleteAccountHandler(repo Repository, logger domain.Logger, ms MailService) *DeleteAccountHandler {
	return &DeleteAccountHandler{repo: repo, logger: logger, ms: ms}
}

// Handle processes the request to delete a user's account.
//
//	@Summary		Delete User Account
//
//	@Description	Delete a user's account
//	@Tags			3- User
//	@Security		BearerAuth
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/users/account [delete]
func (h *DeleteAccountHandler) Handle(ctx context.Context, req *DeleteAccountRequest) (*DeleteAccountResponse, int, error) {

	userId := domain.GetUserID(ctx)

	// deletes user current bucket and  user, also returns email
	fullName, email, err := h.repo.DeleteAccount(ctx, userId)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to delete user account: %v", err))
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	// TODO mail notification
	go func(fullName, email string) {
		maxRetries := 3
		retryInterval := 30 * time.Second

		for attempt := 1; attempt <= maxRetries; attempt++ {

			err := h.ms.SendSuccessfullyDeletedEmail(fullName, email,
				domain.SuccessfullyDeletedEmailSubject,
				domain.EnglishSuccessfullyDeletedEmail)

			if err == nil {
				return
			}
			h.logger.Error("attempt %d failed to send successfully deleted email to %s: %v", attempt, email, err)
			if attempt < maxRetries {
				time.Sleep(retryInterval)
			}
		}
		h.logger.Error("all %d attempts failed for successfully deleted email to %s", maxRetries, email)
		// TODO handle email sending failure (e.g., log it, notify admin, etc.)
	}(fullName, email)

	return nil, http.StatusNoContent, nil
}
