package dto

type Order struct {
	CatalogID uint   `json:"catalog_id" validate:"required,min=1" extensions:"x-order=0"`
	RentDate  string `json:"rent_date" validate:"required" extensions:"x-order=1"`
	Duration  uint   `json:"duration" validate:"required,min=1" extensions:"x-order=2"`
}

type OrderResponse struct {
	CatalogLessDetail               `json:"catalog"`
	RentDate                        string `json:"rent_date" extensions:"x-order=1"`
	ReturnDate                      string `json:"return_date" extensions:"x-order=2"`
	SendInvoiceResponseLessDetailed `json:"invoice"`
}
