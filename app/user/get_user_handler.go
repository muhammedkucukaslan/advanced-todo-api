package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type GetUserRequest struct {
	Id uuid.UUID `params:"id"`
}

type GetUserResponse struct {
	ID              uuid.UUID `json:"id"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	IsEmailVerified bool      `json:"is_email_verified"`
	Address         string    `json:"address"`
}

type GetUserHandler struct {
	repo Repository
}

func NewGetUserHandler(repo Repository) *GetUserHandler {
	return &GetUserHandler{repo: repo}
}

// Handle processes the GetUserRequest and returns the user details.
//
//	@Summary		Get user details by ID for admin
//	@Description	Retrieves user details by ID for admin purposes.
//	@Tags			3- User
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string			true	"User ID"
//	@Success		200	{object}	GetUserResponse	"User details"
//	@Failure		404
//	@Failure		500
//	@Router			/admin/users/{id} [get]
func (h *GetUserHandler) Handle(ctx context.Context, req *GetUserRequest) (*GetUserResponse, int, error) {
	user, err := h.repo.GetUserByIdForAdmin(ctx, req.Id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, 404, err
		}
		return nil, 500, err
	}
	return user, 200, nil
}
