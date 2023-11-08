package controller

import (
	"carstruck/dto"
	"carstruck/entity"
	"carstruck/helpers"
	"carstruck/repository"
	"carstruck/utils"
	"fmt"
	"net/http"
	"strconv"

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
		Address: registerDataTmp.Address,
	}
	
	if err := uc.DbHandler.CreateUser(&registerData); err != nil {
		return err
	}

	verification := entity.Verification{
		UserID: registerData.ID,
		Token: helpers.GenerateVerificationToken(),
	}

	if err := uc.DbHandler.AddToken(&verification); err != nil {
		return err
	}

	if err := helpers.SendVerificationEmail(registerData, verification); err != nil {
		return err
	}
	
	registerDataTmp.Password = ""
	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Registered. Please check your email to do a verification",
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

	if err := uc.DbHandler.CheckVerification(userData); err != nil {
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

func (uc UserController) TopUp(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "Top up",
		Data: user,
	})
}

func (uc UserController) VerifyEmail(c echo.Context) error {
	token := c.Param("token")
	userIdTmp := c.Param("userid")
	userId, err := strconv.Atoi(userIdTmp)
	if err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.Details("Invalid verification URL"))
	}

	verificationData := entity.Verification{
		UserID: uint(userId),
		Token: token,
	}
	
	if err := uc.DbHandler.ValidateEmail(&verificationData); err != nil {
		return err
	}
	
	return c.JSON(http.StatusOK, dto.Response{
		Message: "Validated",
		Data: "Your email has been validated",
	})
}

func (uc UserController) PinpointLocation(c echo.Context) error {
	_, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}
	
	return nil
}