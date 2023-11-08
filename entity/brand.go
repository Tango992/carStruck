package entity

type Brand struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Catalogs []Catalog
}
