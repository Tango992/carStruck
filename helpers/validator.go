package helpers

import (
	"carstruck/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	NewValidator *validator.Validate
}

// Custom validator using go-playground/validator
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.NewValidator.Struct(i); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.Details(err.Error()))
	}
	return nil
}
