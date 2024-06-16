package main

import (
	"chatRoom/routers"
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
	//databaseClient, err := models.DatabaseInit()
	//if err != nil {
	//	log.Fatal(err)
	//}

	e := echo.New()
	// logger
	e.Use(middleware.Logger())

	// adding Websocket and API Routers
	routers.WBRouter(e)
	//routers.APIRouter(e, databaseClient)
	routers.APIRouter(e)

	log.Fatal(e.Start("0.0.0.0:8000"))
}
