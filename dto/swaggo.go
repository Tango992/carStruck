package dto

type GeneralResponse struct {
	Message string `json:"message" extensions:"x-order=0"`
	Data    string `json:"data" extensions:"x-order=1"`
}

type RegisterResponse struct {
	Message      string `json:"message" extensions:"x-order=0"`
	RegisterData `json:"data" extensions:"x-order=1"`
}

type RegisterData struct {
	FullName string `json:"full_name" validate:"required" extensions:"x-order=0"`
	Email    string `json:"email" validate:"email,required" extensions:"x-order=1"`
	Address  string `json:"address" validate:"required" extensions:"x-order=2"`
	Birth    string `json:"birth" validate:"required" extensions:"x-order=3"`
}

type HistoryResponse struct {
	Message string         `json:"message" extensions:"x-order=0"`
	Data    []OrderSummary `json:"data" extensions:"x-order=1"`
}

type CatalogResponse struct {
	Message string    `json:"message" extensions:"x-order=0"`
	Data    []Catalog `json:"data" extensions:"x-order=1"`
}

type OrderResponse struct {
	Message      string `json:"message" extensions:"x-order=0"`
	OrderSummary `json:"data" extensions:"x-order=1"`
}
