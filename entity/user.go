package entity

type User struct {
	ID           uint    `gorm:"primaryKey"`
	FullName     string  `gorm:"not null"`
	Email        string  `gorm:"not null;unique"`
	Password     string  `gorm:"not null"`
	Birth        string  `gorm:"not null;type:date"`
	Address      string  `gorm:"not null"`
	Verification Verification
	Orders       []Order
}
