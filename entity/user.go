package entity

type User struct {
	ID           uint    `gorm:"primaryKey"`
	FullName     string  `gorm:"not null"`
	Email        string  `gorm:"not null;unique"`
	Password     string  `gorm:"not null"`
	Birth        string  `gorm:"not null;type:date"`
	Deposit      float32 `gorm:"not null;default:0"`
	Verification Verification
	Orders       []Order
}
