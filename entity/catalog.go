package entity

type Catalog struct {
	ID         uint    `gorm:"primaryKey"`
	Name       string  `gorm:"not null"`
	Stock      uint    `gorm:"not null;default:0"`
	Cost       float32 `gorm:"not null"`
	CategoryID uint    `gorm:"not null"`
	Histories  []History
}
