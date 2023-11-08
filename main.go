package main

import (
	"carstruck/config"
	"carstruck/entity"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	db := config.InitDb()
	if err := db.AutoMigrate(&entity.User{}, &entity.Validation{}, &entity.Category{}, &entity.Catalog{}, &entity.History{}); err != nil {
		log.Fatal(err)
	}
	
	e := echo.New()

	e.Logger.Fatal(e.Start(":8080"))
}
