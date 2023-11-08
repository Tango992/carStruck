package controller

import (
	"carstruck/dto"
	"carstruck/entity"
	"carstruck/helpers"
	"carstruck/repository"
	"carstruck/utils"
	"fmt"
	"net/http"

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
	var registerDataTmp dto.Register
	if err := c.Bind(&registerDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.Details(err.Error()))
	}

	if err := c.Validate(&registerDataTmp); err != nil {
		return err
	}

	if err := helpers.DateValidator(registerDataTmp.Birth); err != nil {
		return err
	}
	
	if err := helpers.CreateHash(&registerDataTmp); err != nil {
		return err
	}
	
	registerData := entity.User{
		FullName: registerDataTmp.FullName,
		Email: registerDataTmp.Email,
		Password: registerDataTmp.Password,
		Birth: registerDataTmp.Birth,
		Deposit: registerDataTmp.Deposit,
	}

	if err := uc.DbHandler.CreateUser(&registerData); err != nil {
		return err
	}
	
	registerDataTmp.Password = ""
	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Registered",
		Data: registerDataTmp,
	})
}

func (uc UserController) Login(c echo.Context) error {
	var loginData dto.Login
	if err := c.Bind(&loginData); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.Details(err.Error()))
	}

	if err := c.Validate(&loginData); err != nil {
		return err
	}

	userData, err := uc.DbHandler.FindUser(loginData)
	if err != nil {
		return err
	}
	
	if err := helpers.CheckPassword(userData, loginData); err != nil {
		return err
	}

	if err := helpers.SignNewJWT(c, userData); err != nil {
		return err
	}
	
	return c.JSON(http.StatusOK, dto.Response{
		Message: "Logged in",
		Data: fmt.Sprintf("Welcome, %s!", userData.FullName),
	})
}
