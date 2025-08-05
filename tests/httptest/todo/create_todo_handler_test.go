package httptest_todo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	fiberInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/fiber"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	postgresRepo "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTodoHandler(t *testing.T) {

	app := fiber.New()

	tokenService := jwt.NewTokenService(domain.MockJWTTestKey, time.Hour*24, time.Minute*10, time.Minute*10)
	logger := slog.NewLogger()
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

	createTodoHandler := todo.NewCreateTodoHandler(repo)
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
		}, http.StatusBadRequest, domain.ErrInvalidRequest,
		},
		{"empty title", args{
			authHeader: validTokenHeader,
			req:        &todo.CreateTodoRequest{},
		}, http.StatusBadRequest, domain.ErrInvalidRequest,
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
			if tt.code >= 400 {
				assert.NotEmpty(t, resp.Body, "response body should not be empty for error cases")
				var errResp domain.Error
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				require.NoError(t, err)
				assert.NotEmpty(t, errResp.Message, "error message should not be empty")
				assert.Equal(t, errResp.Code, tt.code)
			}
		})
	}
}
