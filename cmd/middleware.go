package main

import (
	"fmt"
	"strings"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"

	"github.com/gofiber/fiber/v2"
)

// Middleware manager to hold shared services
type MiddlewareManager struct {
	tokenService auth.TokenService
}

// NewMiddlewareManager creates a new middleware manager with initialized services
// TODO  Here is dependent on auth.TokenService fix it
func NewMiddlewareManager(tokenService auth.TokenService) *MiddlewareManager {
	return &MiddlewareManager{
		tokenService: tokenService,
	}
}

// AuthMiddleware handles authentication verification
func (m *MiddlewareManager) AuthMiddleware(c *fiber.Ctx) error {
	c.Locals("requireAuth", true)
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Message: "missing authorization header",
			Code:    fiber.StatusUnauthorized,
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Message: "invalid authorization header format",
			Code:    fiber.StatusUnauthorized,
		})
	}

	token := parts[1]
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Message: "missing token",
			Code:    fiber.StatusUnauthorized,
		})
	}

	payload, err := m.tokenService.ValidateToken(token)
	if err != nil {
		fmt.Println("WARN: Invalid or expired token, but allowing anonymous access:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
			Message: "invalid or expired token",
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
			Message: "forbidden resource",
			Code:    fiber.StatusForbidden,
		})
	}
	return c.Next()
}
