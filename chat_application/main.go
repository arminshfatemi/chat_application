package main

import (
	"chatRoom/models"
	"chatRoom/routers/apiRouters"
	"chatRoom/routers/websocketRouters"
	"chatRoom/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// echo setup
	e := echo.New()
	e.Validator = utils.NewValidator()

	// Logger and Recover middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// adding Websocket and API Routers
	apiRouters.AuthAPIRouter(e, databaseClient)
	apiRouters.RoomAPIRouter(e, databaseClient)
	websocketRouters.WBRouter(e)

	log.Fatal(e.Start("0.0.0.0:8000"))
}
