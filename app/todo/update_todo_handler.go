package todo

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type UpdateTodoRequest struct {
	Id    uuid.UUID `params:"id" validate:"required,uuid" swaggerignore:"true"`
	Title string    `json:"title" validate:"required"`
}

type UpdateTodoResponse struct {
}

type UpdateTodoHandler struct {
	repo TodoRepository
}

func NewUpdateTodoHandler(repo TodoRepository) *UpdateTodoHandler {
	return &UpdateTodoHandler{repo: repo}
}

// UpdateTodoHandler handles the update of an existing todo item.
//
//	@Summary		Update an existing todo
//	@Description	Updates an existing todo item for the authenticated user.
//	@Tags			Todo
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id					path	string				true	"Todo ID"
//	@Param			UpdateTodoRequest	body	UpdateTodoRequest	true	"Todo details"
//	@Success		204					"Todo updated successfully"
//	@Failure		400					"Invalid request"
//	@Failure		401					"Unauthorized"
//	@Failure		404					"Todo not found"
//	@Failure		500					"Internal server error"
//	@Router			/todos/{id} [put]
func (h *UpdateTodoHandler) Handle(ctx context.Context, req *UpdateTodoRequest) (*UpdateTodoResponse, int, error) {
	userId := domain.GetUserID(ctx)
	_, err := domain.NewTodo(userId, req.Title)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err = h.repo.UpdateTodo(ctx, req.Id, req.Title); err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusNoContent, nil
}
