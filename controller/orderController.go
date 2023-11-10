package controller

import (
	"carstruck/dto"
	"carstruck/entity"
	"carstruck/helpers"
	"carstruck/repository"
	"carstruck/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	xendit "github.com/xendit/xendit-go/v3"
	invoice "github.com/xendit/xendit-go/v3/invoice"
)

type OrderController struct {
	repository.DbHandler
}

func NewOrderController(dbHandler repository.DbHandler) OrderController {
	return OrderController{
		DbHandler: dbHandler,
	}
}

// Orders  godoc
// @Summary      Submit new car rental order
// @Description  You need an 'Authorization' cookie attached within this request.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      201  {object}  dto.OrderResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      403  {object}  utils.ErrResponse
// @Failure      404  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /orders [post]
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
	returnDate := rentDateFormatted.AddDate(0, 0, int(orderDataTmp.Duration)).Format(dateFormat)

	orderData := entity.Order{
		UserID:     user.ID,
		CatalogID:  orderDataTmp.CatalogID,
		RentDate:   orderDataTmp.RentDate,
		ReturnDate: returnDate,
	}

	subtotal, catalog, err := oc.DbHandler.CreateOrder(&orderData, orderDataTmp.Duration)
	if err != nil {
		return err
	}

	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(fmt.Sprintf("%v", orderData.ID), subtotal)
	xenditClient := xendit.NewClient(os.Getenv("XENDIT_API_KEY"))

	resp, _, errXnd := xenditClient.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()
	if errXnd != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(errXnd.Error()))
	}

	paymentData := entity.Payment{
		OrderID:    orderData.ID,
		InvoiceID:  *resp.Id,
		Amount:     resp.Amount,
		InvoiceURL: resp.InvoiceUrl,
		Status:     resp.Status.String(),
	}

	if err := oc.DbHandler.CreatePayment(&paymentData); err != nil {
		return err
	}

	orderedCatalog := dto.CatalogLessDetail{
		CatalogID: catalog.ID,
		Model:     catalog.Name,
	}

	invoice := dto.SendInvoiceResponseLessDetailed{
		InvoiceId:  paymentData.InvoiceID,
		Amount:     paymentData.Amount,
		InvoiceURL: paymentData.InvoiceURL,
		Status: paymentData.Status,
	}

	orderResponse := dto.OrderSummary{
		OrderID:                         orderData.ID,
		CatalogLessDetail:               orderedCatalog,
		RentDate:                        orderData.RentDate,
		ReturnDate:                      orderData.ReturnDate,
		SendInvoiceResponseLessDetailed: invoice,
	}

	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Order created, proceed to payment",
		Data:    orderResponse,
	})
}

// Orders  godoc
// @Summary      Update payment info from Xendit's server if payment is successful.
// @Tags         orders
// @Accept       json
// @Param        x-callback-token header string true "Secret token from Xendit to validate the request"
// @Param        request body dto.XenditWebhook  true  "Attached data"
// @Success      200
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /orders/update [post]
func (oc OrderController) FetchPaymentUpdate(c echo.Context) error {
	webhookToken := c.Request().Header.Get("x-callback-token")
	if webhookToken != os.Getenv("XENDIT_WEBHOOK_TOKEN") {
		return echo.NewHTTPError(utils.ErrUnauthorized.Details("Invalid webhook token"))
	}
	
	var paymentData dto.XenditWebhook
	if err := c.Bind(&paymentData); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.Details(err.Error()))
	}

	if err := oc.DbHandler.UpdatePaymentStatus(paymentData); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}