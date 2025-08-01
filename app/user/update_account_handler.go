package user

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type UpdateFullNameRequest struct {
	FullName string `json:"full_name,omitempty" validate:"required"`
	Address  string `json:"address,omitempty" validate:"required"`
}

type UpdateFullNameResponse struct{}

type UpdateFullNameService struct {
	repo     Repository
	validate *validator.Validate
}

func NewUpdateFullNameHandler(repo Repository, validate *validator.Validate) *UpdateFullNameService {
	return &UpdateFullNameService{repo: repo, validate: validate}
}

// Handle processes the request to update a user's account information.
//
//	@Summary		Update User Account
//	@Description	Update a user's account information
//	@Tags			3- User
//	@Accept			json
//	@Produce		json
//	@Param			response-language	header	string					true	"Response Language"	Enums(tr, en, ar)
//	@Param			request				body	UpdateFullNameRequest	true	"Update User Account Request"
//	@Security		BearerAuth
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/users/account [patch]
func (h *UpdateFullNameService) Handle(ctx context.Context, req *UpdateFullNameRequest) (*UpdateFullNameResponse, int, error) {
	if err := h.validate.Struct(req); err != nil {
		return nil, 400, domain.ErrInvalidRequest
	}
	userId := domain.GetUserID(ctx)
	if err := h.repo.UpdateFullName(ctx, userId, req.FullName); err != nil {
		fmt.Println("Error updating account:", err)
		return nil, 500, domain.ErrInternalServer
	}
	return nil, 204, nil
}
