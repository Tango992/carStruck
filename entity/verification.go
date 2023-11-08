package entity

type Verification struct {
	UserID    uint   `gorm:"primaryKey"`
	Token     string `gorm:"not null"`
	Validated bool   `gorm:"default:false"`
}
