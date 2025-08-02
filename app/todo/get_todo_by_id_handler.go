package todo

import (
	"context"
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

func NewGetTodoByIdHandler(repo TodoRepository, validator domain.Validator) *GetTodoByIdHandler {
	return &GetTodoByIdHandler{repo: repo}
}

func (h *GetTodoByIdHandler) Handle(ctx context.Context, req *GetTodoByIdRequest) (*GetTodoByIdResponse, int, error) {
	todo, err := h.repo.GetById(ctx, req.Id)
	if err != nil {
		if err == domain.ErrTodoNotFound {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}

	return todo, http.StatusOK, nil
}
