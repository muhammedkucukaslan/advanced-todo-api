//	@title		Advanced Todo API
//	@version	1.0
//	@description
//	@description	## How to use the API
//	@description	1- Click which endpoint you want to use.
//	@description	2- Click "Try it out" button.
//	@description	3- Add your request body or your parameters which are showed and required by the endpoint.
//	@description	4- Click "Execute" button.
//	@description	5- You will see the response.
//	@description
//	@description	Some endpoints require authentication. In this case, you need to log in first.
//	@description	I created two types of users for this project: admin and regular user.
//	@description	Just send a POST request as below at [here](http://localhost:3000/swagger/index.html/).
//	@description	After login, you will get a JWT token in cookies.
//	@description	If you're using cookie-based auth, the cookie will be sent automatically.
//	@description	Alternatively, you can use Bearer Token authentication via the "Authorize" button.
//	@description
//	@description	### Login Request For Admin
//	@description	```json
//	@description	{
//	@description	"email": "admin@admin.com",
//	@description	"password": "admin123"
//	@description	}
//	@description	```
//	@description
//	@description	### Login Request For User
//	@description	```json
//	@description	{
//	@description	"email": "user@user.com",
//	@description	"password": "user1234"
//	@description	}
//	@description	```
//	@description
//	@description	## Error Handling
//	@description	All error responses will follow this JSON format:
//	@description
//	@description	```json
//	@description	{
//	@description	"message": string,
//	@description	"code": int
//	@description	}
//	@description	```
//	@description	### Example
//	@description	```json
//	@description	{
//	@description	"message": "invalid request",
//	@description	"code": 400
//	@description	}
//	@description	```
//	@description
//	@description	## Reminder
//	@description	I did not use `/api` prefix for the endpoint routes. Because I love to host my API on "api" subdomain.
//	@description	Status code with `2xx` is a success code.
//	@description	Status code with `4xx` is a client error code.
//	@description	Status code with `5xx` is a server error code.
//

//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your Bearer token in the format **Bearer &lt;token&gt;**
//
//	@host						localhost:3000

package fiber

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/healthcheck"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/todo"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/user"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	jwtInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	mailersendInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/mailersend"
	postgresInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	redisInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/redis"
	slogInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/slog"
	validatorInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/validator"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	repo, err := postgresInfra.NewRepository(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		panic(err)
	}
	fmt.Println("Connected to database")
	tokenService := jwtInfra.NewTokenService(os.Getenv("JWT_SECRET_KEY"), time.Hour*24, time.Minute*10, time.Minute*10)
	mailersendService := mailersendInfra.NewMailerSendService(os.Getenv("MAILERSEND_API_KEY"), os.Getenv("MAILERSEND_SENDER_EMAIL"), os.Getenv("MAILERSEND_SENDER_NAME"))

	logger := slogInfra.NewLogger()
	validator := validatorInfra.NewValidator(logger)
	redisClient := redisInfra.NewRedisClient(os.Getenv("REDIS_URL"))

	middlewareManager := NewMiddlewareManager(tokenService, logger)

	healthcheckHandler := healthcheck.NewHealthcheckHandler()
	signupHandler := auth.NewSignupHandler(repo, tokenService, mailersendService, validator, logger)
	loginHandler := auth.NewLoginHandler(repo, tokenService, validator, logger)

	getUserHandler := user.NewGetUserHandler(repo)
	getUsersHandler := user.NewGetUsersHandler(repo, validator)
	deleteAccountHandler := user.NewDeleteAccountHandler(repo, logger, mailersendService)
	updateFullNameHandler := user.NewUpdateFullNameHandler(repo, validator)
	getCurrentUserHandler := user.NewGetCurrentUserHandler(repo)
	updatePasswordHandler := user.NewChangePasswordHandler(repo, validator)
	forgotPasswordHandler := user.NewForgotPasswordHandler(repo, mailersendService, tokenService, logger, validator)
	resetPasswordHandler := user.NewResetPasswordHandler(repo, tokenService, logger, validator)
	verifyEmailHandler := user.NewVerifyEmailHandler(repo, validator, tokenService)
	sendVerificationEmailHandler := user.NewSendVerificationEmailHandler(repo, validator, tokenService, mailersendService)

	createTodoHandler := todo.NewCreateTodoHandler(repo, redisClient, logger)
	getTodoByIdHandler := todo.NewGetTodoByIdHandler(repo)
	getTodosHandler := todo.NewGetTodosHandler(repo, redisClient, time.Minute*5)
	updateTodoHandler := todo.NewUpdateTodoHandler(repo)
	deleteTodoHandler := todo.NewDeleteTodoHandler(repo)
	toggleCompletedTodoHandler := todo.NewToggleCompletedTodoHandler(repo)

	app.Get("/healthcheck", Handle(healthcheckHandler, logger))
	app.Use(contextMiddleware)

	adminApp := app.Group("/admin", middlewareManager.AuthMiddleware, middlewareManager.AdminMiddleware)

	app.Post("/signup", Handle(signupHandler, logger))
	app.Post("/login", Handle(loginHandler, logger))

	usersPublicApp := app.Group("/users")
	usersPublicApp.Post("/forgot-password", Handle(forgotPasswordHandler, logger))
	usersPublicApp.Post("/reset-password", Handle(resetPasswordHandler, logger))
	usersPublicApp.Post("/verify-email", Handle(verifyEmailHandler, logger))

	usersApp := app.Group("/users", middlewareManager.AuthMiddleware)
	usersApp.Get("/profile", Handle(getCurrentUserHandler, logger))
	usersApp.Delete("/account", Handle(deleteAccountHandler, logger))
	usersApp.Patch("/account", Handle(updateFullNameHandler, logger))
	usersApp.Patch("/password", Handle(updatePasswordHandler, logger))
	usersApp.Post("/send-verification-email", Handle(sendVerificationEmailHandler, logger))

	usersAdminApp := adminApp.Group("/users")
	usersAdminApp.Get("/", Handle(getUsersHandler, logger))
	usersAdminApp.Get("/:id", Handle(getUserHandler, logger))

	todosApp := app.Group("/todos", middlewareManager.AuthMiddleware)
	todosApp.Post("/", Handle(createTodoHandler, logger))
	todosApp.Get("/:id", Handle(getTodoByIdHandler, logger))
	todosApp.Get("/", Handle(getTodosHandler, logger))
	todosApp.Put("/:id", Handle(updateTodoHandler, logger))
	todosApp.Delete("/:id", Handle(deleteTodoHandler, logger))
	todosApp.Patch("/:id", Handle(toggleCompletedTodoHandler, logger))

	if os.Getenv("ENV") != "production" {
		app.Get("/swagger/*", fiberSwagger.WrapHandler)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Welcome to Advanced Todo API",
			"code":    http.StatusOK,
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Message: "endpoint not found",
			Code:    http.StatusNotFound,
		})
	})
}
