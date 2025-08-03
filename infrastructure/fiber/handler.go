package fiber

import (
	"context"
	"errors"

	"github.com/muhammedkucukaslan/advanced-todo-api/domain"

	"github.com/gofiber/fiber/v2"
)

type Request any
type Response any

type HandlerInterface[R Request, Res Response] interface {
	// int represents the status code
	Handle(ctx context.Context, req *R) (*Res, int, error)
}

func Handle[R Request, Res Response](handler HandlerInterface[R, Res], logger domain.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
			return handleError(c, fiber.StatusBadRequest, err, logger)
		}

		if err := c.CookieParser(&req); err != nil {
			return handleError(c, fiber.StatusBadRequest, err, logger)
		}

		if err := c.ParamsParser(&req); err != nil {
			return handleError(c, fiber.StatusBadRequest, err, logger)
		}

		if err := c.QueryParser(&req); err != nil {
			return handleError(c, fiber.StatusBadRequest, err, logger)
		}

		if err := c.ReqHeaderParser(&req); err != nil {
			return handleError(c, fiber.StatusBadRequest, err, logger)
		}

		if c.Locals("requireAuth") == true {
			role, ok := c.Locals("role").(string)
			if !ok {
				return handleError(c, fiber.StatusUnauthorized, errors.New("invalid role in context"), logger)
			}

			userID, ok := c.Locals("userID").(string)
			if !ok {
				return handleError(c, fiber.StatusUnauthorized, errors.New("invalid user_id in context"), logger)
			}

			ctx := context.WithValue(c.UserContext(), domain.RoleKey, role)
			ctx = context.WithValue(ctx, domain.UserIDKey, userID)
			c.SetUserContext(ctx)
		}

		res, code, err := handler.Handle(c.UserContext(), &req)
		if err != nil {
			return handleError(c, code, err, logger)
		}

		// I send 204 status code if the response is nil
		if res == nil {
			return c.SendStatus(code)
		}

		return c.Status(code).JSON(res)

	}
}

func handleError(c *fiber.Ctx, code int, err error, logger domain.Logger) error {
	var errorType string
	if code >= 500 {
		errorType = "system error"
	} else {
		errorType = "client error"
	}

	logger.Error(err.Error(),
		"type", errorType,
		"method", c.Method(),
		"path", c.Path(),
		"status", code,
		"request_id", c.Locals("requestid"),
	)

	return c.Status(code).JSON(fiber.Map{
		"message": err.Error(),
		"code":    code,
	})

}
