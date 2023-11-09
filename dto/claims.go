package dto

type Claims struct {
	ID       uint    `json:"id"`
	Email    string  `json:"email"`
	FullName string  `json:"full_name"`
	Deposit  float32 `json:"deposit"`
	Address  string  `json:"address"`
}
