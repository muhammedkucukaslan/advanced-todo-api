package todo

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type DeleteTodoRequest struct {
	Id uuid.UUID `params:"id"`
}

type DeleteTodoResponse struct {
}

type DeleteTodoHandler struct {
	repo TodoRepository
}

func NewDeleteTodoHandler(repo TodoRepository) *DeleteTodoHandler {
	return &DeleteTodoHandler{
		repo: repo,
	}
}

// DeleteTodoHandler handles the deletion of a todo item.
//
//	@Summary		Delete a todo
//	@Description	Deletes a todo item for the authenticated user.
//	@Tags			todos
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Todo ID"
//	@Success		204	"Todo deleted successfully"
//	@Failure		401	"Unauthorized"
//	@Failure		500	"Internal server error"
//	@Router			/todos/{id} [delete]
func (h *DeleteTodoHandler) Handle(ctx context.Context, req *DeleteTodoRequest) (*DeleteTodoResponse, int, error) {
	err := h.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return nil, http.StatusNoContent, nil
}
