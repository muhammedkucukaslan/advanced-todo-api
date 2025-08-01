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
//	@description	Just send a POST request as below at [here](http://localhost:3000/swagger/index.html).
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
//	@description	Please handle errors accordingly on the client side.
//	@description	The API returns an error which is according to a language at some endpoints.
//	@description	For example, if you send a request to an anonymous user endpoint, the API will return an error in a specific language.
//	@description	In this case, you need to specify the language in the request header as `accept-language`.
//	@description	I will specify which endpoints require that header.
//	@description
//	@description	If you send a request to an admin endpoint, the API will return an error in Turkish.
//	@description
//	@description				## Reminder
//	@description				I did not use `/api` prefix for the endpoint routes. Because I love to host my API on "api" subdomain.
//	@description				Status code with `2xx` is a success code.
//	@description				Status code with `4xx` is a client error code.
//	@description				Status code with `5xx` is a server error code.
//
//	@securityDefinitions.apikey	JWTAuth
//	@in							cookie
//	@name						jwt
//	@description				JWT cookie obtained from login endpoint
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your Bearer token in the format **Bearer &lt;token&gt;**
//
//	@host						localhost:3000

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/auth"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/healthcheck"
	"github.com/muhammedkucukaslan/advanced-todo-api/app/user"
	"github.com/muhammedkucukaslan/advanced-todo-api/domain"
	fiberInfra "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/fiber"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/jwt"
	mailersend "github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/mailersend"
	"github.com/muhammedkucukaslan/advanced-todo-api/infrastructure/postgres"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	repo, err := postgres.NewRepository(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		panic(err)
	}
	fmt.Println("Connected to database")
	tokenService := jwt.NewTokenService(os.Getenv("JWT_SECRET_KEY"), time.Hour*24, time.Minute*10, time.Minute*10)
	// cookieService := fiberInfra.NewFiberCookieService()
	mailersendService := mailersend.NewMailerSendService(os.Getenv("MAILERSEND_API_KEY"), os.Getenv("MAILERSEND_SENDER_EMAIL"), os.Getenv("MAILERSEND_SENDER_NAME"))
	MockEmailServer := &domain.MockEmailServer{}
	fmt.Println(MockEmailServer)
	validate := validator.New(validator.WithRequiredStructEnabled())
	middlewareManager := NewMiddlewareManager(tokenService)

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	healthcheckHandler := healthcheck.NewHealthcheckHandler()
	signupHandler := auth.NewSignupHandler(repo, tokenService, mailersendService, validate, logger)
	loginHandler := auth.NewLoginHandler(repo, tokenService, validate, logger)

	getUserHandler := user.NewGetUserHandler(repo)
	getUsersHandler := user.NewGetUsersHandler(repo, validate)
	deleteAccountHandler := user.NewDeleteAccountHandler(repo, logger, validate, MockEmailServer)
	updateFullNameHandler := user.NewUpdateFullNameHandler(repo, validate)
	getCurrentUserHandler := user.NewGetCurrentUserHandler(repo)
	updatePasswordHandler := user.NewChangePasswordHandler(repo, validate)
	forgotPasswordHandler := user.NewForgotPasswordHandler(repo, mailersendService, tokenService, logger, validate)
	resetPasswordHandler := user.NewResetPasswordHandler(repo, tokenService, logger, validate)
	verifyEmailHandler := user.NewVerifyEmailHandler(repo, validate, tokenService)
	sendVerificationEmailHandler := user.NewSendVerificationEmailHandler(repo, validate, tokenService, mailersendService)

	app.Get("/healthcheck", handle(healthcheckHandler))
	app.Use(fiberInfra.ContextMiddleware)

	adminApp := app.Group("/admin", middlewareManager.AuthMiddleware, middlewareManager.AdminMiddleware)

	app.Post("/signup", handle(signupHandler))
	app.Post("/login", handle(loginHandler))

	usersPublicApp := app.Group("/users")
	usersPublicApp.Post("/forgot-password", handle(forgotPasswordHandler))
	usersPublicApp.Post("/reset-password", handle(resetPasswordHandler))
	usersPublicApp.Post("/verify-email", handle(verifyEmailHandler))

	usersApp := app.Group("/users", middlewareManager.AuthMiddleware)
	usersApp.Get("/profile", handle(getCurrentUserHandler))
	usersApp.Delete("/account", handle(deleteAccountHandler))
	usersApp.Patch("/account", handle(updateFullNameHandler))
	usersApp.Patch("/password", handle(updatePasswordHandler))
	usersApp.Post("/send-verification-email", handle(sendVerificationEmailHandler))

	usersAdminApp := adminApp.Group("/users")
	usersAdminApp.Get("/", handle(getUsersHandler))
	usersAdminApp.Get("/:id", handle(getUserHandler))

	todosApp := app.Group("/todos")
	todosApp.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Message: "This endpoint is not implemented yet",
			Code:    http.StatusNotImplemented,
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Message: "Welcome to Advanced Todo API",
			Code:    http.StatusOK,
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Message: "endpoint not found",
			Code:    http.StatusNotFound,
		})
	})
}
