package repository

import (
	"carstruck/dto"
	"carstruck/entity"
	"carstruck/utils"
	"errors"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type DbHandler struct {
	*gorm.DB
}

func NewDBHandler(db *gorm.DB) DbHandler {
	return DbHandler{
		DB: db,
	}
}

func (db DbHandler) CreateUser(user *entity.User) error {
	if err := db.Create(user).Error; err != nil {
		return echo.NewHTTPError(utils.ErrConflict.Details(err.Error()))
	}
	return nil
}

func (db DbHandler) FindUser(loginData dto.Login) (entity.User, error) {
	var user entity.User

	res := db.Where("email = ?", loginData.Email).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return entity.User{}, echo.NewHTTPError(utils.ErrUnauthorized.Details("Invalid email / password"))
		}
		return entity.User{}, echo.NewHTTPError(utils.ErrInternalServer.Details(res.Error.Error()))
	}
	return user, nil
}