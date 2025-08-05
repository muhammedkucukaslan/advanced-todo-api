package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type UpdateFullNameRequest struct {
	FullName string `json:"full_name,omitempty" validate:"required"`
	Address  string `json:"address,omitempty" validate:"required"`
}

type UpdateFullNameResponse struct{}

type UpdateFullNameService struct {
	repo     Repository
	validate domain.Validator
}

func NewUpdateFullNameHandler(repo Repository, validate domain.Validator) *UpdateFullNameService {
	return &UpdateFullNameService{repo: repo, validate: validate}
}

// Handle processes the request to update a user's account information.
//
//	@Summary		Update User Account
//	@Description	Update a user's account information
//	@Tags			3- User
//	@Accept			json
//	@Produce		json
//	@Param			request				body	UpdateFullNameRequest	true	"Update User Account Request"
//	@Security		BearerAuth
//	@Success		204
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/users/account [patch]
func (h *UpdateFullNameService) Handle(ctx context.Context, req *UpdateFullNameRequest) (*UpdateFullNameResponse, int, error) {
	if err := h.validate.Validate(req); err != nil {
		return nil, http.StatusBadRequest, domain.ErrInvalidRequest
	}
	userId := domain.GetUserID(ctx)
	if err := h.repo.UpdateFullName(ctx, userId, req.FullName); err != nil {
		fmt.Println("Error updating account:", err)
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}
	return nil, http.StatusNoContent, nil
}
