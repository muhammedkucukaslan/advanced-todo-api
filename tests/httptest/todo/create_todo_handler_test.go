package httptest_todo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	fiberInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/fiber"
	jwtInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	slogInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	testUtils "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/require"
)

func TestCreateTodoHandler(t *testing.T) {

	app := fiber.New()

	tokenService := jwtInfra.NewTokenService(domain.MockJWTTestKey, time.Hour*24, time.Minute*10, time.Minute*10)
	logger := slogInfra.NewLogger()
	middlewareManager := fiberInfra.NewMiddlewareManager(tokenService, logger)
	app.Use(middlewareManager.AuthMiddleware)

	ctx := context.Background()

	postgresContainer, connStr := testUtils.CreateTestContainer(t, ctx)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate postgres container")
	}()

	repo, err := postgresRepo.NewRepository(connStr)
	require.NoError(t, err, "failed to create repository")
	runMigrations(t, connStr)
	setupTestUser(t, connStr)

	// I am not trying to test caching. So, i can use mock.
	createTodoHandler := todo.NewCreateTodoHandler(repo, testUtils.NewMockCache())
	app.Post("/todos", fiberInfra.Handle(createTodoHandler, logger))

	validToken, err := tokenService.GenerateToken(domain.RealUserId, domain.TestUser.Role, time.Now())
	require.NoError(t, err, "failed to generate valid token")

	fakeUserIdToken, err := tokenService.GenerateToken(domain.FakeUserId, domain.TestUser.Role, time.Now())
	require.NoError(t, err, "failed to generate fake token")

	validTokenHeader := "Bearer " + validToken
	fakeUserIdTokenHeader := "Bearer " + fakeUserIdToken

	type args struct {
		authHeader string
		req        *todo.CreateTodoRequest
	}

	tests := []struct {
		name    string
		args    args
		code    int
		wantErr error
	}{
		{
			"valid creation", args{
				authHeader: validTokenHeader,
				req: &todo.CreateTodoRequest{
					Title: "Test Todo",
				},
			}, http.StatusCreated, nil,
		},
		{"fake user ID ", args{
			authHeader: fakeUserIdTokenHeader,
			req: &todo.CreateTodoRequest{
				Title: "Test Todo",
			},
		}, http.StatusNotFound, domain.ErrUserNotFound,
		},
		{"too short title", args{
			authHeader: validTokenHeader,
			req: &todo.CreateTodoRequest{
				Title: "ab",
			},
		}, http.StatusBadRequest, domain.ErrTitleTooShort,
		},
		{"empty title", args{
			authHeader: validTokenHeader,
			req:        &todo.CreateTodoRequest{},
		}, http.StatusBadRequest, domain.ErrEmptyTitle,
		},
		{"too long title", args{
			authHeader: validTokenHeader,
			req: &todo.CreateTodoRequest{
				Title: strings.Repeat("a", 101),
			},
		}, http.StatusBadRequest, domain.ErrTitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.args.req != nil {
				data, err := json.Marshal(tt.args.req)
				require.NoError(t, err)
				body = bytes.NewReader(data)
			}

			req, _ := http.NewRequest(http.MethodPost, "/todos", body)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tt.args.authHeader)

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
