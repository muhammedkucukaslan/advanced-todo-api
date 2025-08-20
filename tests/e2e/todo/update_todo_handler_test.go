package e2etest_todo

import (
	"bytes"
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

type updateTodoHandlerArgs struct {
	authHeader string
	req        *todo.UpdateTodoRequest
}
type updateTodoHandlerTestCase struct {
	name              string
	args              *updateTodoHandlerArgs
	putCode           int
	wantPutErr        error
	getCode           int
	wantGetErr        error
	wantGetResp       *todo.GetTodoByIdResponse
	needsExistingTodo bool
}

func TestUpdateTodoHandler(t *testing.T) {
	app := fiber.New()

	tokenService := testUtils.NewTestJWTTokenService()
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

	updateTodoHandler := todo.NewUpdateTodoHandler(repo)
	getTodoByIdHandler := todo.NewGetTodoByIdHandler(repo)
	app.Put("/todos/:id", fiberInfra.Handle(updateTodoHandler, logger))
	app.Get("/todos/:id", fiberInfra.Handle(getTodoByIdHandler, logger))

	validToken, err := tokenService.GenerateAuthAccessToken(domain.RealUserId, domain.TestUser.Role)
	require.NoError(t, err, "failed to generate valid token")

	validTokenHeader := "Bearer " + validToken

	type args struct {
		authHeader string
		req        *todo.UpdateTodoRequest
	}

	newTitle := "Updated Test Todo"

	// I made id as uuid.Nil.
	// In unchanging test cases, such as "empty title" and "too short title", i am expecting the todo to be unchanged.
	// However, if used a specific UUID, the todo would affected by "valid update" test case while whole process.
	// So, i decided to create a new todo for each test case and then asign the new id to id.
	// You can analyze the process it between line 165 and 173.

	tests := []updateTodoHandlerTestCase{
		{
			name: "valid update",
			args: &updateTodoHandlerArgs{
				authHeader: validTokenHeader,
				req: &todo.UpdateTodoRequest{
					Id:    uuid.Nil,
					Title: newTitle,
				},
			},
			putCode:    http.StatusNoContent,
			wantPutErr: nil,
			getCode:    http.StatusOK,
			wantGetErr: nil,
			wantGetResp: &todo.GetTodoByIdResponse{
				Id:        uuid.Nil,
				Title:     newTitle,
				Completed: domain.TestTodo.Completed,
			},
			needsExistingTodo: true,
		},
		{
			name: "too short title",
			args: &updateTodoHandlerArgs{
				authHeader: validTokenHeader,
				req: &todo.UpdateTodoRequest{
					Id:    uuid.Nil,
					Title: "ab",
				},
			},
			putCode:    http.StatusBadRequest,
			wantPutErr: domain.ErrTitleTooShort,
			getCode:    http.StatusOK,
			wantGetErr: nil,
			wantGetResp: &todo.GetTodoByIdResponse{
				Id:        uuid.Nil,
				Title:     domain.TestTodo.Title,
				Completed: domain.TestTodo.Completed,
			},
			needsExistingTodo: true,
		},
		{
			name: "empty title",
			args: &updateTodoHandlerArgs{
				authHeader: validTokenHeader,
				req: &todo.UpdateTodoRequest{
					Id:    uuid.Nil,
					Title: "",
				},
			},
			putCode:    http.StatusBadRequest,
			wantPutErr: domain.ErrEmptyTitle,
			getCode:    http.StatusOK,
			wantGetErr: nil,
			wantGetResp: &todo.GetTodoByIdResponse{
				Id:        uuid.Nil,
				Title:     domain.TestTodo.Title,
				Completed: domain.TestTodo.Completed,
			},
			needsExistingTodo: true,
		},
		{
			name: "not existed todo",
			args: &updateTodoHandlerArgs{
				authHeader: validTokenHeader,
				req: &todo.UpdateTodoRequest{
					Id:    domain.FakeTodoUuid,
					Title: newTitle,
				},
			},
			putCode:           http.StatusNotFound,
			wantPutErr:        domain.ErrTodoNotFound,
			getCode:           http.StatusNotFound,
			wantGetErr:        domain.ErrTodoNotFound,
			wantGetResp:       nil,
			needsExistingTodo: false,
		},
	}

	for _, tt := range tests {
		todoId := uuid.New()

		// I seperated 2 test cases as "PUT" and "GET" because it became easier to understand which test case is failing.
		t.Run(tt.name+" PUT", func(t *testing.T) {

			if tt.needsExistingTodo {
				setupTestTodoWithTitle(t, todoId, connStr)
				tt.args.req.Id = todoId
				tt.wantGetResp.Id = todoId
			}
			sendTestUpdateRequest(t, app, todoId, &tt)
		})

		t.Run(tt.name+" GET", func(t *testing.T) {
			sendTestGetRequestForUpdatedTodo(t, app, todoId, &tt)
		})
	}
}

func sendTestUpdateRequest(t *testing.T, app *fiber.App, todoId uuid.UUID, tc *updateTodoHandlerTestCase) {
	var body io.Reader
	if tc.args.req != nil {
		data, err := json.Marshal(tc.args.req)
		require.NoError(t, err)
		body = bytes.NewReader(data)
	}

	req := httptest.NewRequest(http.MethodPut, "/todos/"+todoId.String(), body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tc.args.authHeader)

	resp, err := app.Test(req, -1)
	require.NoError(t, err, "failed to send PUT request")
	defer resp.Body.Close()

	require.Equal(t, tc.putCode, resp.StatusCode, "PUT status code mismatch")

	if testUtils.IsErrorStatusCode(tc.putCode) {
		testUtils.VerifyErrorResponse(t, resp.Body, tc.wantPutErr)
	}
}

func sendTestGetRequestForUpdatedTodo(t *testing.T, app *fiber.App, todoId uuid.UUID, tc *updateTodoHandlerTestCase) {
	req := httptest.NewRequest(http.MethodGet, "/todos/"+todoId.String(), nil)
	req.Header.Set("Authorization", tc.args.authHeader)

	resp, err := app.Test(req, -1)
	require.NoError(t, err, "failed to send GET request")
	defer resp.Body.Close()

	require.Equal(t, tc.getCode, resp.StatusCode, "GET status code mismatch")

	if testUtils.IsErrorStatusCode(tc.getCode) {
		testUtils.VerifyErrorResponse(t, resp.Body, tc.wantGetErr)
	} else {
		verifyGetRequestSuccessResponse(t, resp.Body, tc.wantGetResp)
	}
}

func verifyGetRequestSuccessResponse(t *testing.T, body io.ReadCloser, expected *todo.GetTodoByIdResponse) {
	var getResp todo.GetTodoByIdResponse
	err := json.NewDecoder(body).Decode(&getResp)
	require.NoError(t, err, "failed to decode success response")

	assert.Equal(t, expected.Id, getResp.Id, "response id should match")
	assert.Equal(t, expected.Title, getResp.Title, "response title should match")
	assert.Equal(t, expected.Completed, getResp.Completed, "response completed status should match")
}
