package apiHandlers

import (
	"chatRoom/rooms"
	"github.com/labstack/echo/v4"
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
