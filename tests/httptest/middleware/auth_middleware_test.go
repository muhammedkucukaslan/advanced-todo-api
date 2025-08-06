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
	jwtInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	slogInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	testUtils "github.com/muhammedkucukaslan/advanced-todo-api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	app := fiber.New()

	realTokenService := jwtInfra.NewTokenService(domain.MockJWTTestKey, time.Hour*24, time.Minute*10, time.Minute*10)
	fakeTokenService := jwtInfra.NewTokenService("fake-key", time.Hour*24, time.Minute*10, time.Minute*10)

	logger := slogInfra.NewLogger()
	middlewareManager := fiberInfra.NewMiddlewareManager(realTokenService, logger)

	healthCheckHandler := healthcheck.NewHealthcheckHandler()
	app.Get("/healthcheck", middlewareManager.AuthMiddleware, fiberInfra.Handle(healthCheckHandler, logger))

	validToken, err := realTokenService.GenerateToken(domain.RealUserId, domain.TestUser.Role, time.Now())
	require.NoError(t, err, "failed to generate valid token")

	fakeToken, err := fakeTokenService.GenerateToken(domain.FakeUserId, domain.TestUser.Role, time.Now())
	require.NoError(t, err, "failed to generate fake token")

	expiredToken, err := realTokenService.GenerateToken(domain.RealUserId, domain.TestUser.Role, time.Now().Add(-time.Hour*100))
	require.NoError(t, err, "failed to generate expired token")

	validTokenHeader := "Bearer " + validToken
	fakeTokenHeader := "Bearer " + fakeToken
	invalidTokenHeader := "Bearer " + "invalid_token"
	invalidHeader := " Bearerrrrrr " + validToken
	missingTokenHeader := "Bearer "
	expiredTokenHeader := "Bearer " + expiredToken

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
			"valid request", args{
				authHeader: validTokenHeader,
				req:        &healthcheck.HealthcheckRequest{},
			}, &healthcheck.HealthcheckResponse{
				Status: healthcheck.IamAlive,
			}, http.StatusOK, nil,
		},
		{
			"invalid token", args{
				authHeader: invalidTokenHeader,
				req:        &healthcheck.HealthcheckRequest{},
			}, nil, http.StatusUnauthorized, domain.ErrInvalidToken,
		},
		{
			"fake token", args{
				authHeader: fakeTokenHeader,
				req:        &healthcheck.HealthcheckRequest{},
			}, nil, http.StatusUnauthorized, domain.ErrInvalidToken,
		},
		{
			"expired token", args{
				authHeader: expiredTokenHeader,
				req:        &healthcheck.HealthcheckRequest{},
			}, nil, http.StatusUnauthorized, domain.ErrInvalidToken,
		},
		{
			"invalid header", args{
				authHeader: invalidHeader,
				req:        &healthcheck.HealthcheckRequest{},
			}, nil,
			http.StatusUnauthorized, domain.ErrInvalidAuthHeader,
		},
		{
			"missing header", args{
				req: &healthcheck.HealthcheckRequest{},
			}, nil, http.StatusUnauthorized, domain.ErrMissingAuthHeader,
		},
		{
			"missing token", args{
				authHeader: missingTokenHeader,
				req:        &healthcheck.HealthcheckRequest{},
			}, nil, http.StatusUnauthorized, domain.ErrInvalidAuthHeader,
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
			if testUtils.IsErrorStatusCode(tt.code) {
				testUtils.VerifyErrorResponse(t, resp.Body, tt.wantErr)
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
