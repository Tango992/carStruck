package main

import (
	"carstruck/config"
	"carstruck/controller"
	_ "carstruck/docs"
	"carstruck/helpers"
	"carstruck/middlewares"
	"carstruck/repository"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/swaggo/echo-swagger"
)

// @title carStruck API
// @version 1.0
// @description A car rental API utilizing payment gateway (Xendit) and Google Maps Static API. Made as a project for Hacktiv8, derived from myself to give a digitalized business solution.

// @contact.name Daniel Rahmanto
// @contact.email daniel.rahmanto@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host car-struck.fly.dev
// @BasePath /
func main() {
	db := config.InitDb()
	dbHandler := repository.NewDBHandler(db)
	userController := controller.NewUserController(dbHandler)
	orderController := controller.NewOrderController(dbHandler)
	catalogController := controller.NewCatalogController(dbHandler)

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}
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

	e.GET("/", func(c echo.Context) error {return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")})
	
	users := e.Group("/users")
	{
		users.POST("/register", userController.Register)
		users.GET("/verify/:userid/:token", userController.VerifyEmail)
		users.POST("/login", userController.Login)
		users.GET("/pinpoint", userController.PinpointLocation, middlewares.RequireAuth)
		users.GET("/history", userController.History, middlewares.RequireAuth)
	}

	orders := e.Group("/orders")
	{
		orders.POST("", orderController.NewOrder, middlewares.RequireAuth)
		orders.POST("/update", orderController.FetchPaymentUpdate)
	}
	e.GET("/catalogs", catalogController.ViewCatalogHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
