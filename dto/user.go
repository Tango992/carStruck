package dto

type Register struct {
	FullName string  `json:"full_name" validate:"required" extensions:"x-order=0"`
	Email    string  `json:"email" validate:"email,required" extensions:"x-order=1"`
	Password string  `json:"password,omitempty" validate:"required" extensions:"x-order=2"`
	Address  string  `json:"address" validate:"required" extensions:"x-order=3"`
	Birth    string  `json:"birth" validate:"required" extensions:"x-order=4"`
}

type Login struct {
	Email    string `json:"email" validate:"email,required" extensions:"x-order=0"`
	Password string `json:"password,omitempty" validate:"required" extensions:"x-order=1"`
}
