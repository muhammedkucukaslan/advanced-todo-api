package httptest_todo

import (
	"context"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	fiberInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/fiber"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	slogInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	testUtils "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/require"
)

func TestToggleCompletedTodoHandler(t *testing.T) {

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
	setupTestTodo(t, connStr)

	toogleCompletedTodoHandler := todo.NewToggleCompletedTodoHandler(repo)
	app.Patch("/todos/:id", fiberInfra.Handle(toogleCompletedTodoHandler, logger))

	validToken, err := tokenService.GenerateAuthAccessToken(domain.RealUserId, domain.TestUser.Role)
	require.NoError(t, err, "failed to generate valid token")

	validTokenHeader := "Bearer " + validToken

	tests := []struct {
		name    string
		id      string
		code    int
		wantErr error
	}{
		{
			"valid", domain.TestTodo.Id.String(), http.StatusNoContent, nil,
		},
		{
			"todo not found", domain.FakeTodoId, http.StatusNotFound, domain.ErrTodoNotFound,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPatch, "/todos/"+tt.id, nil)

			req.Header.Set("Authorization", validTokenHeader)

			resp, err := app.Test(req, -1)
			require.NoError(t, err, "failed to create request")
			defer resp.Body.Close()

			require.Equal(t, tt.code, resp.StatusCode)
			if testUtils.IsErrorStatusCode(tt.code) {
				testUtils.VerifyErrorResponse(t, resp.Body, tt.wantErr)
			}
		})
	}
}
