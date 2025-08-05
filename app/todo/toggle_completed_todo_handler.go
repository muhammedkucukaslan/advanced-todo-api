package todo

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type ToggleCompletedTodoRequest struct {
	Id uuid.UUID `params:"id"`
}

type ToggleCompletedTodoResponse struct{}

type ToggleCompletedTodoHandler struct {
	repo TodoRepository
}

func NewToggleCompletedTodoHandler(repo TodoRepository) *ToggleCompletedTodoHandler {
	return &ToggleCompletedTodoHandler{repo: repo}
}

// ToggleCompletedTodoHandler handles the toggling of a todo item's completion status.
//
//	@Summary		Toggle todo completion status
//	@Description	Toggles the completion status of a todo item for the authenticated user.
//	@Tags			Todo
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Todo ID"
//
//	@Success		204	"Todo completion status toggled"
//
//	@Failure		400	"Invalid request"
//	@Failure		401	"Unauthorized"
//	@Failure		404	"Todo not found"
//	@Failure		500	"Internal server error"
//	@Router			/todos/{id} [patch]
func (h *ToggleCompletedTodoHandler) Handle(ctx context.Context, req *ToggleCompletedTodoRequest) (*ToggleCompletedTodoResponse, int, error) {
	if err := h.repo.ToggleCompleted(ctx, req.Id); err != nil {
		if err == domain.ErrTodoNotFound {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}
	return nil, http.StatusNoContent, nil
}
