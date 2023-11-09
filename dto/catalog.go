package dto

type Catalog struct {
	CatalogID uint    `json:"catalog_id"`
	Brand     string  `json:"brand"`
	Model     string  `json:"model"`
	Category  string  `json:"category"`
	Stock     uint    `json:"stock"`
	Cost      float32 `json:"cost"`
}
