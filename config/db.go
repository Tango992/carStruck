package config

import (
	"carstruck/entity"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	AutoMigrate(db)
	return db
}

func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&entity.User{}, &entity.Validation{}, &entity.Category{}, &entity.Catalog{}, &entity.History{}); err != nil {
		log.Fatal(err)
	}
}