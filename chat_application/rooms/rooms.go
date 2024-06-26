package rooms

import (
	"chatRoom/models"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	ChatRooms = make(map[string]*ChatRoom)
	Upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Client struct {
	ID   primitive.ObjectID
	conn *websocket.Conn
	room *ChatRoom
	send chan []byte
}

//type MessageJson struct {
//	Type    string `json:"type"`
//	Content string `json:"content"`
//}

type ErrorMessage struct {
	Error string `json:"error"`
}

type ChatRoom struct {
	ID         primitive.ObjectID
	name       string
	clients    map[*Client]bool
	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
}

func CreateNewChatRoom(name string, id primitive.ObjectID) *ChatRoom {
	return &ChatRoom{
		ID:         id,
		name:       name,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run is function that is responsible for controlling the chatRoom like sending message, user joining and leaving chat
// Run will get close if no one is in the room
func (room *ChatRoom) Run(mongoClient *mongo.Client, notificationChannel chan models.Message, redisClient *redis.Client) {
	// delete the Run from the map
	// NOTE: its important because we don't want to have several Run for a single Room
	defer func() {
		delete(ChatRooms, room.name)
	}()

	for {
		select {
		// when a user what to join the chat room
		case client := <-room.register:
			room.clients[client] = true

		// when a user want to leave the room
		case client := <-room.unregister:
			if _, exists := room.clients[client]; exists {
				delete(room.clients, client)
				close(client.send)
			}

			// we will stop runner to reduce load if there is no clients in the room
			if room.ClientsCount() == 0 {
				return
			}

		// case when a new message is sent by users, we save message in database then sent to users
		case message := <-room.broadcast:
			// save the Message in the database
			_, err := mongoClient.Database("chat_app").Collection("messages").InsertOne(context.TODO(), message)
			if err != nil {
				log.Println(err)
				return
			}

			// purge the cash to make new cache
			if err := models.PurgeRecentMessagesCache(redisClient, room.name); err != nil {
				log.Fatalln(err)
			}

			// and send the message to the notification producer channel
			notificationChannel <- message

			for client := range room.clients {
				go WriteMessage(client, message)
			}
		}
	}
}

// ClientsCount will show the count of the clients in the chatRoom
func (room *ChatRoom) ClientsCount() int {
	return len(room.clients)
}

// RegisterClient will register the user to the
func (room *ChatRoom) RegisterClient(client *Client) {
	room.register <- client
}

func CreateNewClient(conn *websocket.Conn, chatRoom *ChatRoom, id primitive.ObjectID) *Client {
	return &Client{
		ID:   id,
		conn: conn,
		room: chatRoom,
		send: make(chan []byte)}
}

// WriteMessage job is to send the message to the other users
func WriteMessage(client *Client, message models.Message) {
	jsonMessage := models.MessageJson{
		Type:    "message",
		Content: message.Content,
	}
	msg, err := json.Marshal(jsonMessage)
	if err != nil {
		log.Println(err)
		return
	}

	err = client.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println(err)
		client.room.unregister <- client
		if err := client.conn.Close(); err != nil {
			log.Println("WriteMessage: ", err)
		}
	}
}

// sendErrorMessage sends an error message back to the client
func sendErrorMessage(client *Client, errorMsg string) {
	errorMessage := ErrorMessage{Error: errorMsg}
	jsonErrorMessage, err := json.Marshal(errorMessage)
	if err != nil {
		log.Println("Failed to marshal error message:", err)
		return
	}

	err = client.conn.WriteMessage(websocket.TextMessage, jsonErrorMessage)
	if err != nil {
		log.Println("Failed to send error message:", err)
	}
}

// ReadMessage main function and goroutine to handle sent messages from clients
func ReadMessage(client *Client) {
	// if user is disconnected we will remove it from the chatroom
	defer func() {
		client.room.unregister <- client
		if err := client.conn.Close(); err != nil {
			log.Println("ReadMessage in defer: ", err)
		}
	}()

	// loop that read the message from the user send it to the
mainLoop:
	for {
		// get the message
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Println(string(message))
			log.Println("ReadMessage in default: ", err)
			break mainLoop
		}

		// unmarshal the Json, send the error if json was invalid,
		var msg models.MessageJson
		err = json.Unmarshal(message, &msg)
		if err != nil {
			sendErrorMessage(client, "Invalid JSON format")
			continue
		}
		messageStruct := models.Message{
			Content:   msg.Content,
			SenderID:  client.ID,
			RoomID:    client.room.ID,
			Timestamp: time.Now(),
		}

		// sending the message to the Room Run
		client.room.broadcast <- messageStruct
	}
}

func GetRecentMessages(client *mongo.Client, roomID primitive.ObjectID) ([]models.MessageJson, error) {
	collection := client.Database("chat_app").Collection("messages")

	filter := bson.M{"room_id": roomID}
	findOptions := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(50)

	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each document into a Message struct
	var messages []models.MessageJson
	if err = cursor.All(context.TODO(), &messages); err != nil {
		return []models.MessageJson{}, nil
	}

	return messages, nil
}

func WriteListMessage(client *Client, messages []models.MessageJson) error {
	msg, err := json.Marshal(messages)
	if err != nil {
		log.Println(err)
		return err
	}

	err = client.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println(err)
		client.room.unregister <- client
		if err := client.conn.Close(); err != nil {
			log.Println("WriteMessage: ", err)
		}
		return err
	}
	return nil
}
