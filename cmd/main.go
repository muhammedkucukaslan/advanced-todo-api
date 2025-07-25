package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	_ "github.com/muhammedkucukaslan/advanced-todo-api/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	env := os.Getenv("ENV")
	if env != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, skipping...")
		}
	}

	PORT := ":3000"

	app := setupServer()
	setupRoutes(app)
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
