package main

import (
	//"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"notification/database"
)

// init function to load environment variables
//func init() {
//	// Load environment variables from .env file
//	if err := godotenv.Load(); err != nil {
//		log.Fatal("Error loading .env file")
//	}
//}

func main() {
	// connecting to the database
	_, err := database.DatabaseInit()
	if err != nil {
		log.Fatal(err)
	}

	// echo setup
	e := echo.New()

	// Logger and Recover middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	log.Fatal(e.Start("0.0.0.0:8000"))
}
