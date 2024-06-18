package apiHandlers

import (
	"chatRoom/models"
	"chatRoom/rooms"
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type CreateRoomRequest struct {
	Name string `json:"name" validate:"required"`
}

// CreateNewRoomHandler will create a ChatRoom to join
func CreateNewRoomHandler(mongoClient *mongo.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var request CreateRoomRequest
		if err := c.Bind(&request); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if err := c.Validate(request); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		// check if there is any room with given name in the database
		var foundRoom models.Room
		clientCollection := mongoClient.Database("chat_app").Collection("room")
		err := clientCollection.FindOne(context.TODO(), bson.M{"name": request.Name}).Decode(&foundRoom)
		if err == nil {
			return c.String(http.StatusBadRequest, "room already exists")
		}

		// TODO if we need this code
		//// check if chatRoom with given name already exists
		//_, exists := rooms.ChatRooms[request.Name]
		//if exists {
		//	return c.String(http.StatusBadRequest, "room already exists")
		//}

		// Insert the room into the database
		databaseRoom := models.CreateNewRoom(request.Name)
		if err := databaseRoom.CreateRoomInDatabase(mongoClient); err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// add the room to the map
		rooms.ChatRooms[request.Name] = rooms.CreateNewChatRoom(request.Name)

		return c.String(http.StatusOK, "room created")
	}
}

// JoinRoomHandler is handler that join the user to the room
func JoinRoomHandler(client *mongo.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomName := c.QueryParam("name")

		// check database to see if room exists
		if err := models.RoomExists(roomName, client); err != nil {
			return c.String(http.StatusBadRequest, "room does not exist")
		}

		// if room is not found then we will add it to the map.
		// if room is not in the map it means it's not Run so we will run it too
		chatRoom, exists := rooms.ChatRooms[roomName]
		if exists != true {
			// add the room to the map
			chatRoom = rooms.CreateNewChatRoom(roomName)
			rooms.ChatRooms[roomName] = chatRoom
			go chatRoom.Run()
		}

		// upgrade the connection and add the client to the chatRoom
		conn, err := rooms.Upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return c.String(http.StatusInternalServerError, "something went wrong")
		}
		client := rooms.CreateNewClient(conn, chatRoom)
		chatRoom.RegisterClient(client)

		// goroutine that listens to messages that are going to be sent by client
		rooms.ReadMessage(client)
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
