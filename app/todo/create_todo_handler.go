package todo

import (
	"context"
	"errors"
	"net/http"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type CreateTodoRequest struct {
	Title string `json:"title" validate:"required,min=1,max=100"`
}

type CreateTodoResponse struct {
}

type CreateTodoHandler struct {
	repo TodoRepository
}

func NewCreateTodoHandler(repo TodoRepository) *CreateTodoHandler {
	return &CreateTodoHandler{repo: repo}
}

// CreateTodoHandler handles the creation of a new todo item.
//
//	@Summary		Create a new todo
//	@Description	Creates a new todo item for the authenticated user.
//	@Tags			todos
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
	if req.Title == "" {
		return nil, http.StatusBadRequest, domain.ErrInvalidRequest
	}

	userId := domain.GetUserID(ctx)

	todo, err := domain.NewTodo(userId, req.Title)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err = h.repo.CreateTodo(ctx, todo); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, http.StatusForbidden, domain.ErrUserNotFound
		}
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusCreated, nil
}
