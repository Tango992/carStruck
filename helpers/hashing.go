package helpers

import (
	"carstruck/dto"
	"carstruck/entity"
	"carstruck/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Create a new password hash
func CreateHash(data *dto.Register) error {
	hashed, err:= bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	data.Password = string(hashed)
	return nil
}

// Returns error if password does not match with the original hash
func CheckPassword(dbData entity.User, data dto.Login) error {
	if err := bcrypt.CompareHashAndPassword([]byte(dbData.Password), []byte(data.Password)); err != nil {
		return echo.NewHTTPError(utils.ErrUnauthorized.Details("Invalid email / password"))
	}
	return nil
}