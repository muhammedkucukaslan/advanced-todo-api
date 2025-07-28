package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func setupServer() *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Concurrency:  256 * 1024,
		AppName:      "Advanced Todo API",
	})

	registerMiddlewares(app)

	if os.Getenv("ENV") == "development" {
		app.Get("/swagger/*", fiberSwagger.WrapHandler)
	}
	return app
}

func registerMiddlewares(app *fiber.App) {

	allowOrigins := os.Getenv("ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:8000"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Content-Type,Authorization,X-Requested-With",
	}))

	app.Use(recover.New())

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${method} | ${path} | ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	app.Use(helmet.New())

	app.Use(limiter.New(limiter.Config{
		Max:        99,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    http.StatusTooManyRequests,
				"message": "you've sent too many requests. Please wait.",
			})
		},
	}))

}
