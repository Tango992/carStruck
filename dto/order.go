package dto

type Order struct {
	CatalogID uint   `json:"catalog_id" validate:"required,min=1" extensions:"x-order=0"`
	RentDate  string `json:"rent_date" validate:"required" extensions:"x-order=1"`
	Duration  uint   `json:"duration" validate:"required,min=1" extensions:"x-order=2"`
}

type OrderSummary struct {
	OrderID                         uint `json:"order_id" extensions:"x-order=0"`
	CatalogLessDetail               `json:"catalog" extensions:"x-order=1"`
	RentDate                        string `json:"rent_date" extensions:"x-order=2"`
	ReturnDate                      string `json:"return_date" extensions:"x-order=3"`
	SendInvoiceResponseLessDetailed `json:"invoice" extensions:"x-order=4"`
}
