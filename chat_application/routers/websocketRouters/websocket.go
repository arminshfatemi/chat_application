package websocketRouters

import (
	"chatRoom/handlers/apiHandlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func WBRouter(e *echo.Echo) {
	r := e.Group("/ws/")
	// protected routes
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		ContextKey: "userToken",
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
	}))

	r.GET("join/", apiHandlers.JoinRoomHandler())
}
