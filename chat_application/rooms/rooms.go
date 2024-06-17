package rooms

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type Client struct {
	conn *websocket.Conn
	room *ChatRoom
	send chan []byte
}

type ChatRoom struct {
	name       string
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func CreateNewChatRoom(name string) *ChatRoom {
	return &ChatRoom{
		name:       name,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}
func WriteMessage(client *Client) {
	defer func() {
		client.room.unregister <- client
		if err := client.conn.Close(); err != nil {
			log.Println("WriteMessage: ", err)
		}
	}()

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				err := client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Println("WriteMessage inside: ", err)
				}
				log.Println("WriteMessage: connection closed")
				return
			}
			client.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// ReadMessage main function and goroutine to handle sent messages from clients
func ReadMessage(client *Client) {
	// if user is disconnected we will remove it from the chatroom
	defer func() {
		client.room.unregister <- client
		if err := client.conn.Close(); err != nil {
			log.Println("ReadMessage: ", err)
		}
	}()

	// TODO: check if need to put a shutdown case in select case
	// loop that read the message from the user send it to the
mainLoop:
	for {
		select {
		default:
			_, message, err := client.conn.ReadMessage()
			if err != nil {
				log.Println("ReadMessage: ", err)
				break mainLoop
			}
			client.room.broadcast <- message
		}
	}
}

// Run is function that is responsible for controlling the chatRoom like sending message, user joining and leaving chat
func (room *ChatRoom) Run() {
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

		// when a user send a message in the room
		case message := <-room.broadcast:
			for client := range room.clients {
				select {
				case client.send <- message:
				default:
					delete(room.clients, client)
					close(client.send)
				}
			}
		}
	}
}

func (room *ChatRoom) ClientsCount() int {
	return len(room.clients)
}

func ServeWs(room *ChatRoom, c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "something went wrong")
	}

	client := &Client{conn: conn, room: room, send: make(chan []byte)}
	room.register <- client
	go ReadMessage(client)
	WriteMessage(client)
	return nil
}
