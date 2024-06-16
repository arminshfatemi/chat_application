package rooms

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

var (
	chatRooms = make(map[string]*ChatRoom)
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func CreateChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Query().Get("name")
	if roomName == "" {
		http.Error(w, "Please provide a room name", http.StatusBadRequest)
		return
	}

	_, exists := chatRooms[roomName]
	if exists {
		http.Error(w, "Room already exists", http.StatusBadRequest)
		return
	}
	chatRoom := CreateNewChatRoom(roomName)
	chatRooms[roomName] = chatRoom

	go chatRoom.Run()
	log.Println("Room created: ", chatRoom.name)
}

func JoinChatRoomHandler(c echo.Context) error {
	roomName := c.Request().URL.Query().Get("name")
	if roomName == "" {
		return c.String(http.StatusBadRequest, "Please provide a room name")
	}

	chatRoom, exists := chatRooms[roomName]
	if !exists {
		return c.String(http.StatusBadRequest, "Room does not exists")
	}
	return ServeWs(chatRoom, c)
}

func ListChatRoomsHandler(w http.ResponseWriter, r *http.Request) {
	clientList := []string{}
	for client, _ := range chatRooms {
		clientList = append(clientList, client)
	}

	_, err := fmt.Fprintln(w, clientList)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
