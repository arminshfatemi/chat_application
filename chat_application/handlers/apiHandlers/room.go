package apiHandlers

import (
	"chatRoom/rooms"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type CreateRoomRequest struct {
	Name string `json:"name" validate:"required"`
}

// CreateNewRoomHandler will create a ChatRoom to join
func CreateNewRoomHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request CreateRoomRequest
		if err := c.Bind(&request); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if err := c.Validate(request); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		// check if chatRoom with given name already exists
		_, exists := rooms.ChatRooms[request.Name]
		if exists {
			return c.String(http.StatusBadRequest, "room already exists")
		}

		rooms.ChatRooms[request.Name] = rooms.CreateNewChatRoom(request.Name)

		return c.String(http.StatusOK, "room created")
	}
}

// JoinRoomHandler is handler that join the user to the room
func JoinRoomHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		//var request JoinRoomRequest
		//if err := c.Bind(&request); err != nil {
		//	return c.String(http.StatusInternalServerError, err.Error())
		//}
		//log.Println("first")
		//if err := c.Validate(request); err != nil {
		//	return c.String(http.StatusBadRequest, err.Error())
		//}

		roomName := c.QueryParam("name")

		log.Println("second")

		// check if there is any room with given name
		chatRoom, exists := rooms.ChatRooms[roomName]
		if exists != true {
			return c.String(http.StatusBadRequest, "room not found")
		}
		log.Println("third")

		// if chatRoom is not run we will run it again
		if chatRoom.IsActive() == false {
			go chatRoom.Run()
			chatRoom.Activator()
		}

		// upgrade the connection and add the client to the chatRoom
		conn, err := rooms.Upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return c.String(http.StatusInternalServerError, "something went wrong")
		}
		client := rooms.CreateNewClient(conn, chatRoom)
		chatRoom.RegisterClient(client)

		// goroutine that listens to messages that are going to be sent by client
		rooms.ReadMessage(client, chatRoom.ContextGiver())
		return c.String(http.StatusNoContent, "")

	}
}

// ListAllChatRoomsHandler shows list of all rooms and its client counts
func ListAllChatRoomsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		roomsList := map[string]int{}
		for key, value := range rooms.ChatRooms {
			roomsList[key] = value.ClientsCount()
		}
		return c.JSON(http.StatusOK, roomsList)
	}

}
