package e2e_todo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	fiberInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/fiber"
	jwtInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	redisInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/redis"
	slogInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	testUtils "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/require"
)

type CacheCreateTest struct {
	authHeader string
	req        *todo.CreateTodoRequest
	code       int
	wantErr    error
}

type CacheGetTest struct {
	authHeader     string
	code           int
	wantErr        error
	wantTodosCount int
}

func TestCreateTodoHandlerCaching(t *testing.T) {
	app := fiber.New()

	tokenService := jwtInfra.NewTokenService(domain.MockJWTTestKey, time.Hour*24, time.Minute*10, time.Minute*10)
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

	redisContainer, redisAddr := testUtils.CreateRedisTestContainer(t, ctx)
	defer func() {
		err := redisContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate redis container")

	}()

	redisClient := redisInfra.NewRedisClient(redisAddr)

	createTodoHandler := todo.NewCreateTodoHandler(repo, redisClient, logger)
	getTodosHandler := todo.NewGetTodosHandler(repo, redisClient, time.Minute*5)
	app.Post("/todos", fiberInfra.Handle(createTodoHandler, logger))
	app.Get("/todos", fiberInfra.Handle(getTodosHandler, logger))

	validToken, err := tokenService.GenerateToken(domain.RealUserId, domain.TestUser.Role, time.Now())
	require.NoError(t, err, "failed to generate valid token")

	validTokenHeader := "Bearer " + validToken

	cacheCreateTest := CacheCreateTest{
		authHeader: validTokenHeader,
		req: &todo.CreateTodoRequest{
			Title: domain.TestTodo.Title,
		},
		code:    http.StatusCreated,
		wantErr: nil,
	}

	addedTodosCount := 5
	createNTodosAndCache(t, repo, redisClient, addedTodosCount)

	CacheGetTest := CacheGetTest{
		authHeader:     validTokenHeader,
		code:           http.StatusOK,
		wantErr:        nil,
		wantTodosCount: addedTodosCount + 1,
	}

	t.Run("Cache Create Test POST", func(t *testing.T) {
		sendTestCreateRequest(t, app, &cacheCreateTest)
	})

	t.Run("Cache Get Test GET", func(t *testing.T) {
		sendTestGetRequestForCreatedTodos(t, app, &CacheGetTest)
	})
}

func sendTestCreateRequest(t *testing.T, app *fiber.App, cct *CacheCreateTest) {
	var body io.Reader
	if cct.req != nil {
		data, err := json.Marshal(cct.req)
		require.NoError(t, err)
		body = bytes.NewReader(data)
	}

	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	req.Header.Set("Authorization", cct.authHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err, "failed to send POST request")
	defer resp.Body.Close()

	require.Equal(t, cct.code, resp.StatusCode, "POST status code mismatch")

	if testUtils.IsErrorStatusCode(cct.code) {
		testUtils.VerifyErrorResponse(t, resp.Body, cct.wantErr)
	}
}

func sendTestGetRequestForCreatedTodos(t *testing.T, app *fiber.App, cgt *CacheGetTest) {
	req := httptest.NewRequest(http.MethodGet, "/todos/", nil)
	req.Header.Set("Authorization", cgt.authHeader)

	resp, err := app.Test(req, -1)
	require.NoError(t, err, "failed to send GET request")
	defer resp.Body.Close()

	require.Equal(t, cgt.code, resp.StatusCode, "GET status code mismatch")

	if testUtils.IsErrorStatusCode(cgt.code) {
		testUtils.VerifyErrorResponse(t, resp.Body, cgt.wantErr)
	}

	var todosResponse todo.GetTodosResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&todosResponse), "failed to decode response body")

	require.Equal(t, cgt.wantTodosCount, len(todosResponse), "GET todos count mismatch")
}

func createNTodosAndCache(t *testing.T, repo todo.TodoRepository, cacheClient domain.Cache, n int) {
	ctx := context.Background()
	var todos []*domain.Todo
	for i := 0; i < n; i++ {
		newTodo, _ := domain.NewTodo(domain.TestUser.Id, domain.TestTodo.Title)
		err := repo.CreateTodo(ctx, newTodo)
		require.NoError(t, err, "failed to create todo")

		todos = append(todos, newTodo)

	}
	data, err := json.Marshal(todos)
	require.NoError(t, err, "failed to marshal todos to JSON")
	err = cacheClient.Set(ctx, "todos:"+domain.TestUser.Id.String(), data, time.Minute*5)
	require.NoError(t, err, "failed to set todos in cache")
}
