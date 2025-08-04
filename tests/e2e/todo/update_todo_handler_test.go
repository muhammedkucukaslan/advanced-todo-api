package e2e_todo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	fiberInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/fiber"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateTodoHandler(t *testing.T) {

	app := fiber.New()

	tokenService := jwt.NewTokenService(domain.MockJWTTestKey, time.Hour*24, time.Minute*10, time.Minute*10)
	logger := slog.NewLogger()
	validator := validator.NewValidator(logger)
	middlewareManager := fiberInfra.NewMiddlewareManager(tokenService, logger)
	app.Use(middlewareManager.AuthMiddleware)

	ctx := context.Background()

	postgresContainer, connStr := createTestContainer(t, ctx)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate postgres container")

	}()
	repo, err := postgresRepo.NewRepository(connStr)
	require.NoError(t, err, "failed to create repository")
	runMigrations(t, connStr)
	setupTestUser(t, connStr)

	updateTodoHandler := todo.NewUpdateTodoHandler(repo, validator)
	getTodoByIdHandler := todo.NewGetTodoByIdHandler(repo)
	app.Put("/todos/:id", fiberInfra.Handle(updateTodoHandler, logger))
	app.Get("/todos/:id", fiberInfra.Handle(getTodoByIdHandler, logger))

	validToken, err := tokenService.GenerateToken(domain.RealUserId, domain.TestUser.Role, time.Now())
	require.NoError(t, err, "failed to generate valid token")

	validTokenHeader := "Bearer " + validToken

	type args struct {
		authHeader string
		req        *todo.UpdateTodoRequest
	}

	newTitle := "Updated Test Todo"
	notExistedTodoTestName := "not existed todo id"

	tests := []struct {
		name        string
		args        args
		getWant     *todo.GetTodoByIdResponse
		putCode     int
		wantPostErr error
		getCode     int
		wantGetErr  error
	}{

		{"valid update", args{
			authHeader: validTokenHeader,
			req: &todo.UpdateTodoRequest{
				Id:    uuid.Nil,
				Title: newTitle,
			},
		}, &todo.GetTodoByIdResponse{
			Id:        uuid.Nil,
			Title:     newTitle,
			Completed: domain.TestTodo.Completed,
		}, http.StatusNoContent, nil, http.StatusOK, nil},
		{"too short title", args{
			authHeader: validTokenHeader,
			req: &todo.UpdateTodoRequest{
				Id:    uuid.Nil,
				Title: "ab",
			},
		}, &todo.GetTodoByIdResponse{
			Id:        uuid.Nil,
			Title:     domain.TestTodo.Title,
			Completed: domain.TestTodo.Completed,
		}, http.StatusBadRequest, domain.ErrTitleTooShort, http.StatusOK, nil},

		{"empty title", args{
			authHeader: validTokenHeader,
			req: &todo.UpdateTodoRequest{
				Id:    uuid.Nil,
				Title: "",
			},
		}, &todo.GetTodoByIdResponse{
			Id:        uuid.Nil,
			Title:     domain.TestTodo.Title,
			Completed: domain.TestTodo.Completed,
		}, http.StatusBadRequest, domain.ErrTitleTooShort, http.StatusOK, nil},

		{notExistedTodoTestName, args{
			authHeader: validTokenHeader,
			req: &todo.UpdateTodoRequest{
				Id:    uuid.Nil,
				Title: newTitle,
			},
		}, nil, http.StatusNotFound, domain.ErrTodoNotFound, http.StatusNotFound, domain.ErrTodoNotFound},
	}

	for _, tt := range tests {
		todoId := uuid.New()

		t.Run(tt.name, func(t *testing.T) {

			if tt.name != notExistedTodoTestName {
				setupTestTodoWithTitle(t, todoId, connStr)
				tt.args.req.Id = todoId
				tt.getWant.Id = todoId
			}

			var body io.Reader
			if tt.args.req != nil {
				data, err := json.Marshal(tt.args.req)
				require.NoError(t, err)
				body = bytes.NewReader(data)
			}

			putReq, _ := http.NewRequest(http.MethodPut, "/todos/"+todoId.String(), body)

			putReq.Header.Set("Content-Type", "application/json")
			putReq.Header.Set("Authorization", tt.args.authHeader)

			resp, err := app.Test(putReq, -1)
			require.NoError(t, err, "failed to create request")

			defer resp.Body.Close()

			require.Equal(t, tt.putCode, resp.StatusCode)
			if tt.putCode >= 400 {
				assert.NotEmpty(t, resp.Body, "response body should not be empty for error cases")
				var errResp domain.Error
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				require.NoError(t, err)
				assert.NotEmpty(t, errResp.Message, "error message should not be empty")
				assert.Equal(t, errResp.Code, tt.putCode)
			}
		})

		t.Run(tt.name+" get", func(t *testing.T) {
			getReq, _ := http.NewRequest(http.MethodGet, "/todos/"+todoId.String(), nil)
			getReq.Header.Set("Authorization", tt.args.authHeader)

			resp, err := app.Test(getReq, -1)
			require.NoError(t, err, "failed to create request")

			defer resp.Body.Close()

			require.Equal(t, tt.getCode, resp.StatusCode)
			if tt.getCode >= 400 {
				assert.NotEmpty(t, resp.Body, "response body should not be empty for error cases")
				var errResp domain.Error
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				require.NoError(t, err)
				assert.NotEmpty(t, errResp.Message, "error message should not be empty")
				assert.Equal(t, errResp.Code, tt.getCode)
			} else {
				var getResp todo.GetTodoByIdResponse
				err = json.NewDecoder(resp.Body).Decode(&getResp)
				require.NoError(t, err)
				assert.Equal(t, tt.getWant.Id, getResp.Id, "response id should match")
				assert.Equal(t, tt.getWant.Title, getResp.Title, "response title should match")
				assert.Equal(t, tt.getWant.Completed, getResp.Completed, "response completed status should match")
			}
		})
	}
}
