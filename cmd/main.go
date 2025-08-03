package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/muhammedkucukaslan/advanced-todo-api/docs"
	fiberInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/fiber"
)

func main() {

	env := os.Getenv("ENV")
	if env != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, skipping...")
		}
	}

	PORT := ":3000"

	app := fiberInfra.SetupServer()
	fiberInfra.SetupRoutes(app)
	fmt.Println("Server is running on port", PORT)

	go func() {
		if err := app.Listen(PORT); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	}()

	gracefulShutdown(app)
}

func gracefulShutdown(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}
}
