package httptest_middleware

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/healthcheck"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	fiberInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/fiber"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminMiddleware(t *testing.T) {

	app := fiber.New()

	realTokenService := jwt.NewTokenService(domain.MockJWTTestKey, time.Hour*24, time.Minute*10, time.Minute*10)

	logger := slog.NewLogger()
	middlewareManager := fiberInfra.NewMiddlewareManager(realTokenService, logger)

	healthCheckHandler := healthcheck.NewHealthcheckHandler()
	app.Get("/healthcheck", middlewareManager.AuthMiddleware, middlewareManager.AdminMiddleware, fiberInfra.Handle(healthCheckHandler, logger))

	adminToken, err := realTokenService.GenerateToken(domain.RealUserId, "ADMIN", time.Now())
	require.NoError(t, err, "failed to generate valid token")

	userToken, err := realTokenService.GenerateToken(domain.RealUserId, "USER", time.Now())
	require.NoError(t, err, "failed to generate valid token")

	type args struct {
		authHeader string
		req        *healthcheck.HealthcheckRequest
	}

	tests := []struct {
		name    string
		args    args
		want    *healthcheck.HealthcheckResponse
		code    int
		wantErr error
	}{
		{
			"valid  request", args{
				authHeader: "Bearer " + adminToken,
				req:        &healthcheck.HealthcheckRequest{},
			}, &healthcheck.HealthcheckResponse{
				Status: healthcheck.IamAlive,
			}, http.StatusOK, nil,
		},
		{
			"invalid request", args{
				authHeader: "Bearer " + userToken,
				req:        &healthcheck.HealthcheckRequest{},
			}, nil, http.StatusForbidden, domain.ErrForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/healthcheck", nil)
			require.NoError(t, err)

			if tt.args.authHeader != "" {
				req.Header.Set("Authorization", tt.args.authHeader)
			}

			resp, err := app.Test(req, -1)
			require.NoError(t, err, "failed to create request")
			require.NotNil(t, resp)
			defer resp.Body.Close()

			require.Equal(t, tt.code, resp.StatusCode)
			if tt.code >= 400 {
				assert.NotEmpty(t, resp.Body, "response body should not be empty for error cases")
				var errResp domain.Error
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				require.NoError(t, err)
				assert.NotEmpty(t, errResp.Message, "error message should not be empty")
				assert.Equal(t, errResp.Message, tt.wantErr.Error())
				assert.Equal(t, errResp.Code, tt.code)
			} else {
				var res healthcheck.HealthcheckResponse
				err = json.NewDecoder(resp.Body).Decode(&res)
				require.NoError(t, err, "failed to decode response")
				assert.NotNil(t, res, "response should not be nil")
				assert.Equal(t, tt.want.Status, res.Status, "response status should match")
			}
		})
	}
}
