package controller

import (
	"carstruck/repository"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	repository.DbHandler
}

func NewUserController(dbHandler repository.DbHandler) UserController {
	return UserController{
		DbHandler: dbHandler,
	}
}

func (uc UserController) Register(c echo.Context) error {
	return nil
}

func (uc UserController) Login(c echo.Context) error {
	return nil
}