package apiRouters

import (
	"chatRoom/handlers/apiHandlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func RoomAPIRouter(e *echo.Echo, mongoClient *mongo.Client) {
	r := e.Group("/api/")

	// protected routes
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		ContextKey: "userToken",
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
	}))

	r.POST("room/create/", apiHandlers.CreateNewRoomHandler())
	r.GET("room/list/", apiHandlers.ListAllChatRoomsHandler())

}
