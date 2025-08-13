package fiber

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
)

type MiddlewareManager struct {
	tokenService auth.TokenService
	logger       domain.Logger
}

// TODO  Here is dependent on auth.TokenService fix it
func NewMiddlewareManager(tokenService auth.TokenService, logger domain.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		tokenService: tokenService,
		logger:       logger,
	}
}

func (m *MiddlewareManager) AuthMiddleware(c *fiber.Ctx) error {
	c.Locals("requireAuth", true)
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Message: domain.ErrMissingAuthHeader.Error(),
			Code:    fiber.StatusUnauthorized,
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Message: domain.ErrInvalidAuthHeader.Error(),
			Code:    fiber.StatusUnauthorized,
		})
	}

	token := parts[1]
	payload, err := m.tokenService.ValidateAuthToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Message: err.Error(),
			Code:    fiber.StatusUnauthorized,
		})
	}
	c.Locals("userID", payload.UserID)
	c.Locals("role", payload.Role)
	return c.Next()
}

func (m *MiddlewareManager) AdminMiddleware(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "ADMIN" {

		return c.Status(fiber.StatusForbidden).JSON(domain.Error{
			Message: domain.ErrForbidden.Error(),
			Code:    fiber.StatusForbidden,
		})
	}
	return c.Next()
}

type FiberContextKey struct{}

func contextMiddleware(c *fiber.Ctx) error {
	ctx := context.WithValue(c.UserContext(), FiberContextKey{}, c)
	c.SetUserContext(ctx)
	return c.Next()
}
