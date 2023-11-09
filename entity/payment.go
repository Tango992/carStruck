package entity

type Payment struct {
	OrderID       uint    `gorm:"primaryKey"`
	InvoiceID     string  `gorm:"not null"`
	Amount        float32 `gorm:"not null"`
	InvoiceURL    string  `gorm:"not null"`
	Status        string  `gorm:"not null"`
	PaymentMethod string
	CreatedAt     string `gorm:"type:timestamptz;not null"`
	CompletedAt   string `gorm:"type:timestamptz;default:null"`
}
