package dto

type Order struct {
	CatalogID uint   `json:"catalog_id" validate:"required,min=1" extensions:"x-order=0"`
	RentDate  string `json:"rent_date" validate:"required" extensions:"x-order=1"`
	Duration  uint   `json:"duration" validate:"required,min=1" extensions:"x-order=2"`
}
