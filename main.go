package main

import (
	"carstruck/config"
	"carstruck/controller"
	"carstruck/helpers"
	"carstruck/repository"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := config.InitDb()
	dbHandler := repository.NewDBHandler(db)
	userController := controller.NewUserController(dbHandler)

	e := echo.New()
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	users := e.Group("/users")
	{
		users.POST("/register", userController.Register)
		users.POST("/login", userController.Login)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
