package main

import (
	"chatRoom/models"
	"chatRoom/routers"
	"chatRoom/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"log"
)

// init function to load environment variables
func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// connecting to the database
	databaseClient, err := models.DatabaseInit()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	// Register the custom validator to validate the clients request body
	e.Validator = utils.NewValidator()

	// Logger and Recover middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// adding Websocket and API Routers
	routers.WBRouter(e)
	routers.APIRouter(e, databaseClient)

	log.Fatal(e.Start("0.0.0.0:8000"))
}
