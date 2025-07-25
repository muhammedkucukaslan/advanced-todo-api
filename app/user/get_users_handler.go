package user

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type GetUsersRequest struct {
	Page  int `query:"page" validate:"required,min=1"`
	Limit int `query:"limit" validate:"required,min=1,max=100"`
}

type GetUsersResponse []User

type User struct {
	Id       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
}

type GetUsersHandler struct {
	repo     Repository
	validate *validator.Validate
}

func NewGetUsersHandler(repo Repository, validate *validator.Validate) *GetUsersHandler {
	return &GetUsersHandler{
		repo:     repo,
		validate: validate,
	}
}

// Handle processes the GetUsersRequest and returns a list of users.
//
//	@Summary		Get users for admin
//	@Description	Fetch users with pagination
//	@Tags			3- User
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			page	query		int	true	"Page number"
//	@Param			limit	query		int	true	"Page size"
//	@Success		200		{object}	GetUsersResponse
//	@Failure		400
//	@Failure		500
//	@Router			/admin/users [get]
func (h *GetUsersHandler) Handle(ctx context.Context, req *GetUsersRequest) (*GetUsersResponse, int, error) {

	if err := h.validate.Struct(req); err != nil {
		return nil, http.StatusBadRequest, domain.ErrInvalidRequest
	}

	users, err := h.repo.GetUsers(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, http.StatusInternalServerError, domain.ErrInternalServer
	}

	return &users, http.StatusOK, nil
}
