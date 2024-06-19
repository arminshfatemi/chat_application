package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func JoinNotificationHandler(mongoClient *mongo.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// upgrade the connection
		conn, err := Upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return c.String(http.StatusInternalServerError, "something went wrong")
		}
		defer func() {
			if err := conn.Close(); err != nil {
				log.Println(err)
			}
		}()

		return c.String(200, "done")
	}
}
