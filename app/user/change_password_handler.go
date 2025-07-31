package user

import (
	"context"
	"net/http"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type ChangePasswordResponse struct{}

type ChangePasswordHandler struct {
	repo     Repository
	validate domain.Validator
}

func NewChangePasswordHandler(repo Repository, validate domain.Validator) *ChangePasswordHandler {
	return &ChangePasswordHandler{repo: repo, validate: validate}
}

// Handle processes the request to change a user's password.
//
//	@Summary		Change User Password
//	@Description	Change the password of a user
//	@Tags			3- User
//	@Accept			json
//	@Produce		json
//	@Param			response-language	header	string					true	"Response Language"	enums(tr,ar,en)
//	@Param			request				body	ChangePasswordRequest	true	"Change User Password Request"
//	@Security		BearerAuth
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		404
//	@Failure		500
//	@Router			/users/password [patch]
func (h *ChangePasswordHandler) Handle(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, int, error) {
	userId := domain.GetUserID(ctx)

	if err := h.validate.Validate(req); err != nil {
		return nil, http.StatusBadRequest, domain.ErrInvalidRequest
	}

	if len(req.OldPassword) < 8 || len(req.NewPassword) < 8 {
		return nil, http.StatusBadRequest, domain.ErrPasswordTooShort
	}

	// return a user only having password
	user, err := h.repo.GetUserOnlyHavingPasswordById(ctx, userId)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, http.StatusNotFound, domain.ErrUserNotFound
		}
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	if err := user.ValidatePassword(req.OldPassword); err != nil {
		return nil, http.StatusBadRequest, domain.ErrInvalidCredentials
	}

	if err := user.HashPassword(req.NewPassword); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	user.Id = userId

	if err := h.repo.ChangePassword(ctx, user); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusNoContent, nil
}
