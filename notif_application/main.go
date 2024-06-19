package main

import (
	//"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"notification/database"
	"notification/message_broker"
	"notification/routers"
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
	mongoClient, err := database.ConnectingDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// echo setup
	e := echo.New()

	// Logger and Recover middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Start the Consumer
	go message_broker.RabbitMQConsumer()

	routers.AuthAPIRouter(e, mongoClient)

	log.Fatal(e.Start("0.0.0.0:8080"))
}
