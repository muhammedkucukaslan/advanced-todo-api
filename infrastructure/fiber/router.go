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
	postgresRepo := postgresInfra.NewRepository(os.Getenv("DATABASE_URL"))
	fmt.Println("Connected to database")

	tokenServiceConfig := jwtInfra.Config{
		AccessTokenSecretKey:      os.Getenv("JWT_ACCESS_TOKEN_SECRET_KEY"),
		RefreshTokenSecretKey:     os.Getenv("JWT_REFRESH_TOKEN_SECRET_KEY"),
		AuthAccessTokenDuration:   time.Minute * 15,
		AuthRefreshTokenDuration:  time.Hour * 24 * 30,
		EmailVerificationDuration: time.Minute * 11,
		ForgotPasswordDuration:    time.Minute * 11,
	}

	jwtTokenService := jwtInfra.NewJWTTokenService(tokenServiceConfig)

	mailersendService := mailersendInfra.NewMailerSendService(os.Getenv("MAILERSEND_API_KEY"), os.Getenv("MAILERSEND_SENDER_EMAIL"), os.Getenv("MAILERSEND_SENDER_NAME"))

	slogLogger := slogInfra.NewLogger()
	sl := slogLogger
	validator := validatorInfra.NewValidator(slogLogger)
	redisClient := redisInfra.NewRedisClient(os.Getenv("REDIS_URL"))

	middlewareManager := NewMiddlewareManager(jwtTokenService, slogLogger)

	healthcheckHandler := healthcheck.NewHealthcheckHandler()

	if !domain.IsProdEnv() {
		domain.CookieSecure = false
	}

	signupHandler := auth.NewSignupHandler(&auth.SignupConfig{
		Repo:          postgresRepo,
		TokenService:  jwtTokenService,
		CookieService: NewCookieService(),
		EmailService:  mailersendService,
		Validator:     validator,
		Logger:        slogLogger,
	})

	loginHandler := auth.NewLoginHandler(&auth.LoginConfig{
		Repo:          postgresRepo,
		TokenService:  jwtTokenService,
		CookieService: NewCookieService(),
		Validator:     validator,
		Logger:        slogLogger,
	})

	logoutHandler := auth.NewLogoutHandler(postgresRepo, NewCookieService())
	refreshTokenHandler := auth.NewRefreshTokenHandler(postgresRepo, jwtTokenService)
	getUserHandler := user.NewGetUserHandler(postgresRepo)
	getUsersHandler := user.NewGetUsersHandler(postgresRepo, validator)
	deleteAccountHandler := user.NewDeleteAccountHandler(postgresRepo, sl, mailersendService)
	updateFullNameHandler := user.NewUpdateFullNameHandler(postgresRepo, validator)
	getCurrentUserHandler := user.NewGetCurrentUserHandler(postgresRepo)
	updatePasswordHandler := user.NewChangePasswordHandler(postgresRepo, validator)
	forgotPasswordHandler := user.NewForgotPasswordHandler(postgresRepo, mailersendService, jwtTokenService, sl, validator)
	resetPasswordHandler := user.NewResetPasswordHandler(postgresRepo, jwtTokenService, sl, validator)
	verifyEmailHandler := user.NewVerifyEmailHandler(postgresRepo, validator, jwtTokenService)
	sendVerificationEmailHandler := user.NewSendVerificationEmailHandler(postgresRepo, validator, jwtTokenService, mailersendService)

	createTodoHandler := todo.NewCreateTodoHandler(postgresRepo, redisClient, sl)
	getTodoByIdHandler := todo.NewGetTodoByIdHandler(postgresRepo)
	getTodosHandler := todo.NewGetTodosHandler(postgresRepo, redisClient, time.Minute*5)
	updateTodoHandler := todo.NewUpdateTodoHandler(postgresRepo)
	deleteTodoHandler := todo.NewDeleteTodoHandler(postgresRepo)
	toggleCompletedTodoHandler := todo.NewToggleCompletedTodoHandler(postgresRepo)

	app.Get("/healthcheck", Handle(healthcheckHandler, sl))
	app.Use(contextMiddleware)

	adminApp := app.Group("/admin", middlewareManager.AuthMiddleware, middlewareManager.AdminMiddleware)

	authApp := app.Group("/auth")
	authApp.Post("/signup", Handle(signupHandler, sl))
	authApp.Post("/login", Handle(loginHandler, sl))
	authApp.Post("/logout", Handle(logoutHandler, sl))
	authApp.Post("/refresh", Handle(refreshTokenHandler, sl))

	usersPublicApp := app.Group("/users")
	usersPublicApp.Post("/forgot-password", Handle(forgotPasswordHandler, sl))
	usersPublicApp.Post("/reset-password", Handle(resetPasswordHandler, sl))
	usersPublicApp.Post("/verify-email", Handle(verifyEmailHandler, sl))

	usersApp := app.Group("/users", middlewareManager.AuthMiddleware)
	usersApp.Get("/profile", Handle(getCurrentUserHandler, sl))
	usersApp.Delete("/account", Handle(deleteAccountHandler, sl))
	usersApp.Patch("/account", Handle(updateFullNameHandler, sl))
	usersApp.Patch("/password", Handle(updatePasswordHandler, sl))
	usersApp.Post("/send-verification-email", Handle(sendVerificationEmailHandler, sl))

	usersAdminApp := adminApp.Group("/users")
	usersAdminApp.Get("/", Handle(getUsersHandler, sl))
	usersAdminApp.Get("/:id", Handle(getUserHandler, sl))

	todosApp := app.Group("/todos", middlewareManager.AuthMiddleware)
	todosApp.Post("/", Handle(createTodoHandler, sl))
	todosApp.Get("/:id", Handle(getTodoByIdHandler, sl))
	todosApp.Get("/", Handle(getTodosHandler, sl))
	todosApp.Put("/:id", Handle(updateTodoHandler, sl))
	todosApp.Delete("/:id", Handle(deleteTodoHandler, sl))
	todosApp.Patch("/:id", Handle(toggleCompletedTodoHandler, sl))

	if !domain.IsProdEnv() {
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
