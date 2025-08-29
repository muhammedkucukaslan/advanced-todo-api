package e2etest_todo

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	fiberInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/fiber"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	slogInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	testUtils "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type deleteTodoHandlerArgs struct {
	authHeader string
}
type deleteTodoHandlerTestCase struct {
	name              string
	args              *deleteTodoHandlerArgs
	deleteCode        int
	wantDeleteErr     error
	getCode           int
	wantGetErr        error
	wantGetResp       *todo.GetTodoByIdResponse
	needsExistingTodo bool
}

func TestDeleteTodoHandler(t *testing.T) {
	app := fiber.New()

	tokenService := testUtils.NewTestJWETokenService()
	logger := slogInfra.NewLogger()
	middlewareManager := fiberInfra.NewMiddlewareManager(tokenService, logger)
	app.Use(middlewareManager.AuthMiddleware)

	ctx := context.Background()

	postgresContainer, connStr := testUtils.CreatePostgresTestContainer(t, ctx)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate postgres container")

	}()

	repo := postgresRepo.NewRepository(connStr)
	runMigrations(t, connStr)
	setupTestUser(t, connStr)

	deleteTodoHandler := todo.NewDeleteTodoHandler(repo)
	getTodoByIdHandler := todo.NewGetTodoByIdHandler(repo)
	app.Delete("/todos/:id", fiberInfra.Handle(deleteTodoHandler, logger))
	app.Get("/todos/:id", fiberInfra.Handle(getTodoByIdHandler, logger))

	validToken, err := tokenService.GenerateAuthAccessToken(domain.RealUserId, domain.TestUser.Role)
	require.NoError(t, err, "failed to generate valid token")

	validTokenHeader := "Bearer " + validToken

	tests := []deleteTodoHandlerTestCase{
		{
			name: "valid delete",
			args: &deleteTodoHandlerArgs{
				authHeader: validTokenHeader,
			},
			deleteCode:        http.StatusNoContent,
			wantDeleteErr:     nil,
			getCode:           http.StatusNotFound,
			wantGetErr:        domain.ErrTodoNotFound,
			wantGetResp:       nil,
			needsExistingTodo: true,
		},
		{
			name: "not existed todo",
			args: &deleteTodoHandlerArgs{
				authHeader: validTokenHeader,
			},
			deleteCode:        http.StatusNoContent,
			wantDeleteErr:     nil,
			getCode:           http.StatusNotFound,
			wantGetErr:        domain.ErrTodoNotFound,
			wantGetResp:       nil,
			needsExistingTodo: false,
		},
	}

	for _, tt := range tests {
		todoId := uuid.New()

		t.Run(tt.name+" DELETE", func(t *testing.T) {

			if tt.needsExistingTodo {
				setupTestTodoWithTitle(t, todoId, connStr)
			}
			sendTestDeleteRequest(t, app, todoId, &tt)
		})

		t.Run(tt.name+" GET", func(t *testing.T) {
			sendTestGetRequestForDeletedTodo(t, app, todoId, &tt)
		})
	}
}

func sendTestDeleteRequest(t *testing.T, app *fiber.App, todoId uuid.UUID, tc *deleteTodoHandlerTestCase) {
	req := httptest.NewRequest(http.MethodDelete, "/todos/"+todoId.String(), nil)
	req.Header.Set("Authorization", tc.args.authHeader)

	resp, err := app.Test(req, -1)
	require.NoError(t, err, "failed to send DELETE request")
	defer resp.Body.Close()

	require.Equal(t, tc.deleteCode, resp.StatusCode, "DELETE status code mismatch")

	if testUtils.IsErrorStatusCode(tc.deleteCode) {
		testUtils.VerifyErrorResponse(t, resp.Body, tc.wantDeleteErr)
	}
}

func sendTestGetRequestForDeletedTodo(t *testing.T, app *fiber.App, todoId uuid.UUID, tc *deleteTodoHandlerTestCase) {
	req := httptest.NewRequest(http.MethodGet, "/todos/"+todoId.String(), nil)
	req.Header.Set("Authorization", tc.args.authHeader)

	resp, err := app.Test(req, -1)
	require.NoError(t, err, "failed to send GET request")
	defer resp.Body.Close()

	require.Equal(t, tc.getCode, resp.StatusCode, "GET status code mismatch")

	if testUtils.IsErrorStatusCode(tc.getCode) {
		testUtils.VerifyErrorResponse(t, resp.Body, tc.wantGetErr)
	} else {
		verifyGetRequestSuccessResponseForDeletedTodo(t, resp.Body, tc.wantGetResp)
	}
}

func verifyGetRequestSuccessResponseForDeletedTodo(t *testing.T, body io.ReadCloser, expected *todo.GetTodoByIdResponse) {
	var getResp todo.GetTodoByIdResponse
	err := json.NewDecoder(body).Decode(&getResp)
	require.NoError(t, err, "failed to decode success response")

	assert.Equal(t, expected.Id, getResp.Id, "response id should match")
	assert.Equal(t, expected.Title, getResp.Title, "response title should match")
	assert.Equal(t, expected.Completed, getResp.Completed, "response completed status should match")
}
