package entity

type Order struct {
	ID         uint   `gorm:"not null"`
	UserID     uint   `gorm:"not null"`
	CatalogID  uint   `gorm:"not null"`
	RentDate   string `gorm:"type:date;not null"`
	ReturnDate string `gorm:"type:date;not null"`
	Finished   bool   `gorm:"default:false"`
}
