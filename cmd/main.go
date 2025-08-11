package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/muhammedkucukaslan/advanced-todo-api/docs"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	fiberInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/fiber"
)

func main() {

	if domain.IsProdEnv() {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, skipping...")
		}
	}

	app := fiberInfra.SetupServer()
	fiberInfra.SetupRoutes(app)

	go startServer(app)

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

func startServer(app *fiber.App) {
	PORT := ":3000"

	if err := app.Listen(PORT); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v\n", err))
	}

	fmt.Println("Server is running at http://localhost:", PORT)
}
