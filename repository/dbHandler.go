package repository

import (
	"carstruck/dto"
	"carstruck/entity"
	"carstruck/utils"
	"errors"
	"strconv"
	"time"

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
	if err := db.Table("verifications v").
		Select("v.validated").
		Where("u.id = ?", user.ID).
		Joins("JOIN users u ON v.user_id = u.id").
		Take(&validated).
	Error; err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}

	if !validated {
		return echo.NewHTTPError(utils.ErrForbidden.Details("Please do an email verification first"))
	}
	return nil
}

func (db DbHandler) CreateOrder(data *entity.Order, duration uint) (float32, entity.Catalog, error){
	catalog := entity.Catalog{ID: data.CatalogID}
	
	txErr := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&catalog).Error; err != nil {
			return echo.NewHTTPError(utils.ErrNotFound.Details("Catalog ID does not exist"))
		}
		
		if err := tx.Model(&catalog).Update("stock", gorm.Expr("stock - 1")).Error; err != nil {
			return echo.NewHTTPError(utils.ErrForbidden.Details("Stock is empty"))
		}
		
		if err := tx.Create(data).Error; err != nil {
			return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
		}
		return nil
	})
	if txErr != nil {
		return 0, entity.Catalog{}, txErr
	}

	subtotal := catalog.Cost * float32(duration)
	return subtotal, catalog, nil
}

func (db DbHandler) CreatePayment(data *entity.Payment) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	data.CreatedAt = now
	
	if err := db.Create(data).Error; err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	return nil
}

func (db DbHandler) FindAllCatalogs() ([]dto.Catalog, error) {
	catalogs := []dto.Catalog{}

	if err := db.Table("catalogs c").
		Select("c.id AS catalog_id, b.name AS brand, c.name AS model, cr.name as category, c.stock, c.cost").
		Joins("JOIN categories cr ON cr.id = c.category_id").
		Joins("JOIN brands b ON b.id = c.brand_id").
		Scan(&catalogs).
	Error; err != nil {
		return []dto.Catalog{}, echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	return catalogs, nil
}

func (db DbHandler) FindCatalogByBrand(brand string)  ([]dto.Catalog, error) {
	catalogs := []dto.Catalog{}

	res := db.Table("catalogs c").
		Where("b.name = ?", brand).
		Select("c.id AS catalog_id, b.name AS brand, c.name AS model, cr.name as category, c.stock, c.cost").
		Joins("JOIN categories cr ON cr.id = c.category_id").
		Joins("JOIN brands b ON b.id = c.brand_id").
		Scan(&catalogs)
	err := res.Error

	if err != nil {
		return []dto.Catalog{}, echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}

	if res.RowsAffected == 0 {
		return []dto.Catalog{}, echo.NewHTTPError(utils.ErrNotFound.Details("Brand does not exist in our catalog"))
	}
	
	return catalogs, nil
}

func (db DbHandler) FindCatalogByModel(model string)  ([]dto.Catalog, error) {
	catalog := []dto.Catalog{}

	res := db.
		Table("catalogs c").
		Where("c.name = ?", model).
		Select("c.id AS catalog_id, b.name AS brand, c.name AS model, cr.name as category, c.stock, c.cost").
		Joins("JOIN categories cr ON cr.id = c.category_id").
		Joins("JOIN brands b ON b.id = c.brand_id").
		Take(&catalog)
	err := res.Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return []dto.Catalog{}, echo.NewHTTPError(utils.ErrNotFound.Details("Model does not exist in our catalog"))
		}
		return []dto.Catalog{}, echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	return catalog, nil
}

func (db DbHandler) FindCatalogByBrandAndModel(brand, model string) ([]dto.Catalog, error) {
	catalog := []dto.Catalog{}

	res := db.Table("catalogs c").
		Where("c.name = ?", model).
		Where("b.name = ?", brand).
		Select("c.id AS catalog_id, b.name AS brand, c.name AS model, cr.name as category, c.stock, c.cost").
		Joins("JOIN categories cr ON cr.id = c.category_id").
		Joins("JOIN brands b ON b.id = c.brand_id").
		Take(&catalog)
	err := res.Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return []dto.Catalog{}, echo.NewHTTPError(utils.ErrNotFound.Details("Model does not exist in our catalog"))
		}
		return []dto.Catalog{}, echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	return catalog, nil
}

func (db DbHandler) FindUserOrderHistory(userId uint) ([]dto.OrderSummary, error) {
	orders := []dto.OrderSummary{}

	if err := db.Table("orders o").
		Where("o.user_id = ?", userId).
		Select("o.id AS order_id, c.id AS catalog_id, c.name AS model, o.rent_date, o.return_date, p.invoice_id, p.amount, p.invoice_url, p.status").
		Joins("JOIN catalogs c ON c.id = o.catalog_id").
		Joins("JOIN payments p ON p.order_id = o.id").
		Order("order_id DESC").
		Scan(&orders).
	Error; err != nil {
		return []dto.OrderSummary{}, echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	return orders, nil
}

func (db DbHandler) UpdatePaymentStatus(data dto.XenditWebhook) (error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	orderId, _ := strconv.Atoi(data.ExternalId)
	payment := entity.Payment{OrderID: uint(orderId)}

	if err := db.Model(&payment).
		Updates(&entity.Payment{
			Status: data.Status,
			PaymentMethod: data.PaymentMethod,
			CompletedAt: now,
		}).
	Error; err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.Details(err.Error()))
	}
	return nil
}