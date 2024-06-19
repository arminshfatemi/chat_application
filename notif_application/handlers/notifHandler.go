package handlers

import (
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"notification/database"
)

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func JoinNotificationHandler(mongoClient *mongo.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get the id of user from the claims of JWT token
		token, ok := c.Get("userToken").(*jwt.Token)
		if !ok {
			return c.String(http.StatusUnauthorized, "JWT token missing or invalid")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.String(http.StatusUnauthorized, "failed to cast claims as jwt.MapClaims")
		}
		userId, ok := claims["id"].(string)

		userObjectID, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return c.String(http.StatusInternalServerError, "error in Object id formatting")
		}

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

		client := database.CreateNewClient(conn, userObjectID)
		database.ConnectedClients[userObjectID.Hex()] = client
		client.Run(mongoClient)

		return c.String(200, "done")
	}
}
