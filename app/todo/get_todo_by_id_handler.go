package todo

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type GetTodoByIdRequest struct {
	Id uuid.UUID `params:"id" validate:"required,uuid"`
}

type GetTodoByIdResponse struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
}

type GetTodoByIdHandler struct {
	repo TodoRepository
}

func NewGetTodoByIdHandler(repo TodoRepository) *GetTodoByIdHandler {
	return &GetTodoByIdHandler{repo: repo}
}

// GetTodoByIdHandler handles the retrieval of a todo item by its ID.
//
//	@Summary		Get a todo by ID
//	@Description	Retrieves a todo item by its ID for the authenticated user.
//	@Tags			Todo
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Todo ID"
//	@Success		200	{object}	GetTodoByIdResponse
//	@Failure		400	"Invalid request"
//	@Failure		401	"Unauthorized"
//	@Failure		404	"Todo not found"
//	@Failure		500	"Internal server error"
//	@Router			/todos/{id} [get]
func (h *GetTodoByIdHandler) Handle(ctx context.Context, req *GetTodoByIdRequest) (*GetTodoByIdResponse, int, error) {
	todo, err := h.repo.GetById(ctx, req.Id)
	if err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}

	return todo, http.StatusOK, nil
}
