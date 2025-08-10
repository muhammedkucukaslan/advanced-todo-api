package todo

import (
	"context"
	"encoding/json"
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
	repo  TodoRepository
	cache domain.Cache
	ttl   time.Duration
}

func NewGetTodosHandler(repo TodoRepository, cache domain.Cache, ttl time.Duration) *GetTodosHandler {
	return &GetTodosHandler{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
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
	cacheKey := domain.NewTodoCacheKey(userID)

	if cached, err := h.cache.Get(ctx, cacheKey); err == nil && !isCacheEmpty(cached) {
		var todos GetTodosResponse
		if err := json.Unmarshal(cached, &todos); err == nil {
			return &todos, http.StatusOK, nil
		}
	}

	todos, err := h.repo.GetTodosByUserID(ctx, userID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	go h.SetCache(cacheKey, todos)

	return todos, http.StatusOK, nil
}

func (h *GetTodosHandler) SetCache(key string, todos *GetTodosResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if data, err := json.Marshal(todos); err == nil {
		if err := h.cache.Set(ctx, key, data, h.ttl); err != nil {
			// TODO log cache set error
		}
	}
}

func isCacheEmpty(cached []byte) bool {
	return len(cached) == 0
}
