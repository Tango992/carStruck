package controller

import (
	"carstruck/dto"
	"carstruck/entity"
	"carstruck/helpers"
	"carstruck/repository"
	"carstruck/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type OrderController struct {
	repository.DbHandler
}

func NewOrderController(dbHandler repository.DbHandler) OrderController {
	return OrderController{
		DbHandler: dbHandler,
	}
}

func (oc OrderController) NewOrder(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	var orderDataTmp dto.Order
	if err := c.Bind(&orderDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.Details(err.Error()))
	}

	if err := c.Validate(&orderDataTmp); err != nil {
		return err
	}

	if err := helpers.DateValidator(orderDataTmp.RentDate); err != nil {
		return err
	}
	
	dateFormat := "2006-01-02"
	rentDateFormatted, _ := time.Parse(dateFormat, orderDataTmp.RentDate)
	currentDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Unix()

	if currentDate > rentDateFormatted.Unix() {
		return echo.NewHTTPError(utils.ErrBadRequest.Details("Rent date cannot be less than today"))
	}
	returnDate := time.Now().AddDate(0, 0, int(orderDataTmp.Duration)).Format(dateFormat)
	
	orderData := entity.Order{
		UserID: user.ID,
		CatalogID: orderDataTmp.CatalogID,
		RentDate: orderDataTmp.RentDate,
		ReturnDate: returnDate,
	}

	if err := oc.DbHandler.CreateOrder(&orderData); err != nil {
		return err
	}
	
	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Order created, proceed to payment",
		Data: orderData,
	})
}