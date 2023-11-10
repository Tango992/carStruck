package main_test

import (
	"carstruck/config"
	"carstruck/controller"
	"carstruck/helpers"
	"carstruck/repository"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	e                     *echo.Echo
	db                    *gorm.DB
	dbHandlerTest         repository.DbHandler
	userControllerTest    controller.UserController
	orderControllerTest   controller.OrderController
	catalogControllerTest controller.CatalogController
)

var (
	registerData        = `{"full_name":"Jane","email":"daniel.rahmanto@gmail.com","password":"secret","birth":"2001-01-01","address":"The Breeze BSD"}`
	invalidCredential   = `{"email":"abcdefg@mail.com","password":"abcdefghijklmno"}`
	validCredential     = `{"email":"john@mail.com","password":"secret"}`
	loginResponse       = `{"message":"Logged in","data":"Welcome, John Doe!"}` + "\n"
	verifyEmailResponse = `{"message":"Validated","data":"Your email has been validated"}` + "\n"
	catalogsResponse    = `{"message":"View catalog","data":[{"catalog_id":1,"brand":"Toyota","model":"Avanza","category":"MPV","stock":100,"cost":500000}]}` + "\n"
	updatePayment       = `{"id":"ABCD","external_id":"1","payment_method":"BANK_TRANSFER","status":"PAID"}`
)

func TestMain(m *testing.M) {
	dsn := os.Getenv("DB_URL_TESTING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal(err)
	}
	config.AutoMigrate(db)

	e = echo.New()
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}

	dbHandlerTest = repository.NewDBHandler(db)
	userControllerTest = controller.NewUserController(dbHandlerTest)
	orderControllerTest = controller.NewOrderController(dbHandlerTest)
	catalogControllerTest = controller.NewCatalogController(dbHandlerTest)

	m.Run()
}

func TestRegister(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "localhost:8080/users/register", strings.NewReader(registerData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, userControllerTest.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestEmailVerification(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "localhost:8080/users/verify/:userid/:token", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userid", "token")
	c.SetParamValues("1", "abcd")

	if assert.NoError(t, userControllerTest.VerifyEmail(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, verifyEmailResponse, rec.Body.String())
	}

	c2 := e.NewContext(req, rec)
	c2.SetParamNames("userid", "token")
	c2.SetParamValues("5", "efgh")
	if assert.Error(t, userControllerTest.VerifyEmail(c2)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, verifyEmailResponse, rec.Body.String())
	}
}

func TestLogin(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "localhost:8080/users/login", strings.NewReader(validCredential))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, userControllerTest.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, loginResponse, rec.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodPost, "localhost:8080/users/login", strings.NewReader(invalidCredential))
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c2 := e.NewContext(req2, rec)

	if assert.Error(t, userControllerTest.Login(c2)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, loginResponse, rec.Body.String())
	}
}


func TestGetCatalog(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "localhost:8080/catalogs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, catalogControllerTest.ViewCatalogHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, catalogsResponse, rec.Body.String())
	}
}

func TestXenditWebhook(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "localhost:8080/orders/update", strings.NewReader(updatePayment))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, orderControllerTest.FetchPaymentUpdate(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
