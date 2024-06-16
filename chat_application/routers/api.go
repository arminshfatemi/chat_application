package routers

import (
	"github.com/labstack/echo/v4"
)

//	func APIRouter(e *echo.Echo, mongoClient *mongo.Client) {
//		// user authentication
//		e.POST("user/signup/", func(c echo.Context) error {
//			return c.String(200, "signup")
//		})
//
//		//http.HandleFunc("/create", rooms.CreateChatRoomHandler)
//		//http.HandleFunc("/list", rooms.ListChatRoomsHandler)
//	}
func APIRouter(e *echo.Echo) {
	// user authentication
	e.POST("user/signup/", func(c echo.Context) error {
		return c.String(200, "signup")
	})

	//http.HandleFunc("/create", rooms.CreateChatRoomHandler)
	//http.HandleFunc("/list", rooms.ListChatRoomsHandler)
}
