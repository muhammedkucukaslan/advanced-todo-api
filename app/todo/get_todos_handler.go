package todo

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type GetTodosRequest struct {
}

type GetTodosResponse []Todo

type Todo struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
}

type GetTodosHandler struct {
	repo TodoRepository
}

func NewGetTodosHandler(repo TodoRepository) *GetTodosHandler {
	return &GetTodosHandler{
		repo: repo,
	}
}

// Handle retrieves all todos for the authenticated user.
//
//	@Summary		Get all todos
//	@Description	Retrieves all todos for the authenticated user.
//	@Tags			Todo
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	GetTodosResponse
//	@Failure		401	"Unauthorized"
//	@Failure		500	"Internal server error"
//	@Router			/todos [get]
func (h *GetTodosHandler) Handle(ctx context.Context, req *GetTodosRequest) (*GetTodosResponse, int, error) {
	userID := domain.GetUserID(ctx)
	todos, err := h.repo.GetTodosByUserID(ctx, userID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return todos, http.StatusOK, nil
}
