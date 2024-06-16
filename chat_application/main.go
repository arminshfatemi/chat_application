package main

import (
	"chatRoom/routers"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	e := echo.New()
	// logger
	e.Use(middleware.Logger())

	// adding Websocket and API Routers
	routers.WBRouter(e)
	routers.APIRouter(e)

	log.Fatal(e.Start("127.0.0.1:8000"))
}
