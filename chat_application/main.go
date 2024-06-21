package main

import (
	"chatRoom/broker"
	"chatRoom/models"
	"chatRoom/routers/apiRouters"
	"chatRoom/routers/websocketRouters"
	"chatRoom/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

var (
	notificationProducerChan = make(chan models.Message)
)

func main() {
	// connecting to the database
	databaseClient, err := models.DatabaseInit()
	if err != nil {
		log.Fatal(err)
	}

	// start the redis
	_ = models.InitRedis()

	// start the notification producer
	go broker.NotificationProducer(notificationProducerChan)

	// echo setup
	e := echo.New()
	e.Validator = utils.NewValidator()

	// Logger and Recover middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// adding Websocket and API Routers
	apiRouters.AuthAPIRouter(e, databaseClient)
	apiRouters.RoomAPIRouter(e, databaseClient)
	websocketRouters.WBRouter(e, databaseClient, notificationProducerChan)

	log.Fatal(e.Start("0.0.0.0:8000"))
}
