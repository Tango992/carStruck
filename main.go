package main

import (
	"carstruck/config"
    "github.com/labstack/echo/v4"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config.InitDb()

	e := echo.New()

	e.Logger.Fatal(e.Start(":8080"))
}
