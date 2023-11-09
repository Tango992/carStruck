package controller

import (
	"carstruck/dto"
	"carstruck/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CatalogController struct {
	repository.DbHandler
}

func NewCatalogController(dbHandler repository.DbHandler) CatalogController {
	return CatalogController{
		DbHandler: dbHandler,
	}
}

// View Catalog  godoc
// @Summary      Get catalogs
// @Tags         catalogs
// @Produce      json
// @Param        brand  query  string  false  "Search by brand"
// @Param        model  query  string  false  "Search by model"
// @Success      200  {object}  dto.CatalogResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      404  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /catalogs [get]
func (cc CatalogController) ViewCatalogHandler(c echo.Context) error {
	toTitle := cases.Title(language.AmericanEnglish)
	brand := toTitle.String(c.QueryParam("brand"))
	model := toTitle.String(c.QueryParam("model"))

	if brand != "" && model == "" {
		return cc.ViewByBrand(c, brand)
	}

	if brand == "" && model != "" {
		return cc.ViewByModel(c, model)
	}
	
	if brand != "" && model != "" {
		return cc.ViewSpecific(c, brand, model)
	}
	
	catalogs, err := cc.DbHandler.FindAllCatalogs()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "View catalog",
		Data: catalogs,
	})
}

func (cc CatalogController) ViewSpecific(c echo.Context, brand, model string) error {
	catalogSpecific, err := cc.DbHandler.FindCatalogByBrandAndModel(brand, model)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "View specific catalog",
		Data: catalogSpecific,
	})
}

func (cc CatalogController) ViewByModel(c echo.Context, model string) error {
	catalogsByModel, err := cc.DbHandler.FindCatalogByModel(model)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "View catalog by model",
		Data: catalogsByModel,
	})
}

func (cc CatalogController) ViewByBrand(c echo.Context, brand string) error {
	catalogsByBrand, err := cc.DbHandler.FindCatalogByBrand(brand)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "View catalog by brand",
		Data: catalogsByBrand,
	})
}
