package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"notification/handlers"
	"os"
)

func WBRouter(e *echo.Echo, mongoClient *mongo.Client) {
	r := e.Group("/ws/")

	// protected routes
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		ContextKey: "userToken",
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
	}))

	r.GET("join-notification/", handlers.JoinNotificationHandler(mongoClient))

}
