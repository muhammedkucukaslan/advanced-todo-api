package user

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type ChangePasswordRequest struct {
	Language    string `reqHeader:"response-language" validate:"required,oneof=tr en ar" swaggerignore:"true"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type ChangePasswordResponse struct{}

type ChangePasswordHandler struct {
	repo     Repository
	validate *validator.Validate
}

func NewChangePasswordHandler(repo Repository, validate *validator.Validate) *ChangePasswordHandler {
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

	if err := h.validate.Struct(req); err != nil {
		return nil, 400, domain.ErrInvalidRequest
	}

	if len(req.OldPassword) < 8 || len(req.NewPassword) < 8 {
		return nil, 400, domain.ErrPasswordTooShort
	}

	// return a user only having password
	user, err := h.repo.GetUserOnlyHavingPasswordById(ctx, userId)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, 404, domain.ErrUserNotFound
		}
		return nil, 500, domain.ErrInternalServer
	}
	if err := user.ValidatePassword(req.OldPassword); err != nil {
		return nil, 400, domain.ErrInvalidCredentials
	}

	if err := user.HashPassword(req.NewPassword); err != nil {
		return nil, 500, err
	}

	user.Id = userId

	if err := h.repo.ChangePassword(ctx, user); err != nil {
		return nil, 500, err
	}

	return nil, 204, nil
}
