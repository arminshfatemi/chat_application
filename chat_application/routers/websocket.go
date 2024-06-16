package routers

import (
	"chatRoom/rooms"
	"github.com/labstack/echo/v4"
)

func WBRouter(e *echo.Echo) {
	e.GET("join/", rooms.JoinChatRoomHandler)

}
