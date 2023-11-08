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

func (db DbHandler) AddToken(data *entity.Verification) error {
	if err := db.Create(data).Error; err != nil {
		return echo.NewHTTPError(utils.ErrConflict.Details(err.Error()))
	}
	return nil
}

func (db DbHandler) ValidateEmail(data *entity.Verification) error {
	txErr := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("token = ?", data.Token).First(data).Error; err != nil {
			return echo.NewHTTPError(utils.ErrBadRequest.Details(err.Error()))
		}

		if err := tx.Model(data).Update("validated", true).Error; err != nil {
			return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
		}
		
		return nil
	})
	if txErr != nil {
		return txErr
	}
	return nil
}

func (db DbHandler) CheckVerification(user entity.User) error {
	var validated bool
	if err := db.Table("verifications v").Select("v.validated").Where("u.id = ?", user.ID).Joins("JOIN users u ON v.user_id = u.id").Take(&validated).Error; err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}

	if !validated {
		return echo.NewHTTPError(utils.ErrForbidden.Details("Please do an email verification first"))
	}
	return nil
}

func (db DbHandler) CreateOrder(data *entity.Order) error {
	if err := db.Create(data).Error; err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	return nil
}