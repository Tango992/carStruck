package main

import (
	"carstruck/config"
	"carstruck/controller"
	"carstruck/helpers"
	"carstruck/middlewares"
	"carstruck/repository"
	"os"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func main() {
	db := config.InitDb()
	dbHandler := repository.NewDBHandler(db)
	userController := controller.NewUserController(dbHandler)

	e := echo.New()
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}

	e.Pre(middleware.RemoveTrailingSlash())
	
	logger := zerolog.New(os.Stdout)

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Str("Method", v.Method).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))	
	e.Use(middleware.Recover())

	users := e.Group("/users")
	{
		users.POST("/register", userController.Register)
		users.POST("/login", userController.Login)
		users.POST("/topup", userController.TopUp, middlewares.RequireAuth)
		users.GET("/verify/:userid/:token", userController.VerifyEmail)
	}

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
