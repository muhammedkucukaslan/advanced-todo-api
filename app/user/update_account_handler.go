package user

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type UpdateAccountRequest struct {
	Language string `reqHeader:"response-language" validate:"required,oneof=tr en ar" swaggerignore:"true"`
	FullName string `json:"full_name,omitempty" validate:"required"`
	Address  string `json:"address,omitempty" validate:"required"`
	Phone    string `json:"phone,omitempty" validate:"required"`
}

type UpdateAccountResponse struct{}

type UpdateAccountHandler struct {
	repo     Repository
	validate *validator.Validate
}

func NewUpdateAccountHandler(repo Repository, validate *validator.Validate) *UpdateAccountHandler {
	return &UpdateAccountHandler{repo: repo, validate: validate}
}

// Handle processes the request to update a user's account information.
//
//	@Summary		Update User Account
//	@Description	Update a user's account information
//	@Tags			3- User
//	@Accept			json
//	@Produce		json
//	@Param			response-language	header	string					true	"Response Language"	Enums(tr, en, ar)
//	@Param			request				body	UpdateAccountRequest	true	"Update User Account Request"
//	@Security		BearerAuth
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/users/account [put]
func (h *UpdateAccountHandler) Handle(ctx context.Context, req *UpdateAccountRequest) (*UpdateAccountResponse, int, error) {
	if err := h.validate.Struct(req); err != nil {
		return nil, 400, domain.ErrInvalidRequest
	}
	userID, _ := domain.GetUserID(ctx)
	if err := h.repo.UpdateAccount(ctx, userID, req.FullName, req.Address, req.Phone); err != nil {
		fmt.Println("Error updating account:", err)
		return nil, 500, domain.ErrInternalServer
	}
	return nil, 204, nil
}
