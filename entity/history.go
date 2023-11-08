package entity

import "time"

type History struct {
	ID        uint `gorm:"not null"`
	UserID    uint `gorm:"not null"`
	CatalogID uint `gorm:"not null"`
	Duration  uint `gorm:"not null"`
	CreatedAt time.Time 
}
