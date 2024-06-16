package routers

import (
	"chatRoom/handlers/apiHandlers"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func APIRouter(e *echo.Echo, mongoClient *mongo.Client) {
	// user authentication
	e.POST("user/signup/", apiHandlers.ClientSignUpHandler(mongoClient))

	//http.HandleFunc("/create", rooms.CreateChatRoomHandler)
	//http.HandleFunc("/list", rooms.ListChatRoomsHandler)
}

//http.HandleFunc("/create", rooms.CreateChatRoomHandler)
//http.HandleFunc("/list", rooms.ListChatRoomsHandler)
