package entity

type Payment struct {
	OrderID     uint `gorm:"primaryKey"`
	InvoiceID   uint
	Amount      float32 `gorm:"not null"`
	Status      string
	CreatedAt   string `gorm:"type:timestamp;not null"`
	CompletedAt string `gorm:"type:timestamp"`
}
