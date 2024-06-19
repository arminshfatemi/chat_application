package websocketRouters

import (
	"chatRoom/handlers/apiHandlers"
	"chatRoom/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func WBRouter(e *echo.Echo, mongoDB *mongo.Client, notificationChan chan models.Message) {
	r := e.Group("/ws/")
	// protected routes
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		ContextKey: "userToken",
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
	}))

	r.GET("join/", apiHandlers.JoinRoomHandler(mongoDB, notificationChan))
}
