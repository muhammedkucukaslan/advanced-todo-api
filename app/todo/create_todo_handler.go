package todo

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type CreateTodoRequest struct {
	Title string `json:"title" validate:"required,min=1,max=100"`
}

type CreateTodoResponse struct {
}

type CreateTodoHandler struct {
	repo   TodoRepository
	cache  domain.Cache
	logger domain.Logger
}

func NewCreateTodoHandler(repo TodoRepository, cache domain.Cache, logger domain.Logger) *CreateTodoHandler {
	return &CreateTodoHandler{repo: repo, cache: cache, logger: logger}
}

// CreateTodoHandler handles the creation of a new todo item.
//
//	@Summary		Create a new todo
//	@Description	Creates a new todo item for the authenticated user.
//	@Tags			Todo
//
//	@Security		BearerAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			CreateTodoRequest	body	CreateTodoRequest	true	"Todo details"
//	@Success		201					"Todo created successfully"
//	@Failure		400					"Invalid request"
//	@Failure		401					"Unauthorized"
//	@Failure		500					"Internal server error"
//	@Router			/todos [post]
func (h *CreateTodoHandler) Handle(ctx context.Context, req *CreateTodoRequest) (*CreateTodoResponse, int, error) {
	userId := domain.GetUserID(ctx)

	todo, err := domain.NewTodo(userId, req.Title)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err = h.repo.CreateTodo(ctx, todo); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, http.StatusNotFound, domain.ErrUserNotFound
		}
		return nil, http.StatusInternalServerError, err
	}

	go h.DeleteCacheKey(userId)

	return nil, http.StatusCreated, nil
}

func (h *CreateTodoHandler) DeleteCacheKey(userId uuid.UUID) {
	if err := h.cache.Delete(context.Background(), "todos:"+userId.String()); err != nil {
		h.logger.Error("failed to delete cache key", "key", "todos:"+userId.String(), "error", err)
	}
}
