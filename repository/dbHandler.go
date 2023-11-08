package repository

import "gorm.io/gorm"

type DbHandler struct {
	*gorm.DB
}

func NewDBHandler(db *gorm.DB) DbHandler {
	return DbHandler{
		DB: db,
	}
}