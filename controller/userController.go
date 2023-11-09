package controller

import (
	"carstruck/dto"
	"carstruck/entity"
	"carstruck/helpers"
	"carstruck/repository"
	"carstruck/utils"
	"fmt"
	"io"
	"net/http"
	"os"
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

// Register      godoc
// @Summary      Register new user into database
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.Register  true  "Register data"
// @Success      201  {object}  dto.RegisterResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      409  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /users/register [post]
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

// Login         godoc
// @Summary      Log in with existing account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.Login  true  "Login data"
// @Success      200  {object}  dto.GeneralResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      403  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /users/login [post]
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

// Verification  godoc
// @Summary      Verify newly registered email
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User Id"
// @Param        token   path      int  true  "Secret Token"
// @Success      200  {object}  dto.GeneralResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /users/verify/{userid}/{token} [get]
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

// Get Histories  godoc
// @Summary      Get user's order history
// @Description  You need an 'Authorization' cookie attached within this request.
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.HistoryResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /users/history [get]
func (uc UserController) History(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	orderHistories, err := uc.DbHandler.FindUserOrderHistory(user.ID)
	if err != nil {
		return err
	}
	
	return c.JSON(http.StatusOK, dto.Response{
		Message: "Order History",
		Data: orderHistories,
	})
}

// Get User's location     godoc
// @Summary      Get user's map location
// @Description  You need an 'Authorization' cookie attached within this request.
// @Tags         users
// @Accept       json
// @Produce      image/png
// @Success      200  
// @Failure      500  {object}  utils.ErrResponse
// @Router       /users/pinpoint [get]
func (uc UserController) PinpointLocation(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}
	
	baseUrl := "https://maps.googleapis.com/maps/api/staticmap"
	req, _ := http.NewRequest(http.MethodGet, baseUrl, nil)

	q := req.URL.Query()
	q.Add("center", user.Address)
	q.Add("key", os.Getenv("MAPS_API_KEY"))
	q.Add("markers", "|" + user.Address)
	q.Add("size", "640x640")
	q.Add("zoom", "16")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	defer res.Body.Close()
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	return c.Blob(http.StatusOK, "image/png", body)
}