package controller

import (
	// "bytes"
	"bytes"
	"carstruck/dto"
	"carstruck/entity"
	"carstruck/helpers"
	"carstruck/repository"
	"carstruck/utils"

	// "context"
	"encoding/json"
	"fmt"

	// "fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	// xendit "github.com/xendit/xendit-go/v3"
	// invoice "github.com/xendit/xendit-go/v3/invoice"
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

	subtotal, err := oc.DbHandler.CreateOrder(&orderData, orderDataTmp.Duration)
	if err != nil {
		return err
	}
	
	sendInvoice := dto.SendInvoice{
		ExternalID: fmt.Sprintf("%v", orderData.ID),
		Amount: subtotal,
		Description: user.FullName + "'s order",
		CustomerDetail: dto.CustomerDetail{
			GivenNames: user.FullName,
			Email: user.Email,
		},
	}
	jsonBody, err := json.Marshal(sendInvoice)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	bodyReader := bytes.NewReader(jsonBody)
	
	url := "https://api.xendit.co/v2/invoices"
	base64ApiKey := helpers.BasicAuth64(os.Getenv("XENDIT_API_KEY"), "")
	
	req, _ := http.NewRequest("POST", url, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", base64ApiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	defer res.Body.Close()

	var xenditResponse dto.SendInvoiceResponse
	if err := json.NewDecoder(res.Body).Decode(&xenditResponse); err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	
	// if err := oc.DbHandler.CreatePayment(orderData.ID); err != nil {
	// 	return err
	// }

	// createInvoiceRequest := *invoice.NewCreateInvoiceRequest("test", subtotal) // [REQUIRED] | CreateInvoiceRequest
    // xenditClient := xendit.NewClient("XENDIT_API_KEY"+":")

    // resp, r, err := xenditClient.InvoiceApi.CreateInvoice(context.Background()).
    //     CreateInvoiceRequest(createInvoiceRequest).
    //     Execute()

	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", err.Error())

	// 	b, _ := json.Marshal(err.Error())
	// 	fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

	// 	fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	// }
	
	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Order created, proceed to payment",
		Data: xenditResponse,
	})
}