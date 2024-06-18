package rooms

import (
	"github.com/gorilla/websocket"
	"log"
)

var (
	ChatRooms = make(map[string]*ChatRoom)
	Upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
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

// AddClient add the given client in the clients map
//func (room *ChatRoom) AddClient(name string) {
//	room.clients[]
//}

// WriteMessage job is to send the message to the other users
func WriteMessage(client *Client, message []byte) {
	err := client.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println(err)
		client.room.unregister <- client
		if err := client.conn.Close(); err != nil {
			log.Println("WriteMessage: ", err)
		}
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
		select {
		default:
			_, message, err := client.conn.ReadMessage()
			if err != nil {
				log.Println(string(message))
				log.Println("ReadMessage in default: ", err)
				break mainLoop
			}
			log.Println("after error handling")
			client.room.broadcast <- message
		}
	}
}

// Run is function that is responsible for controlling the chatRoom like sending message, user joining and leaving chat
func (room *ChatRoom) Run() {
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

		// case when a new message is sent by users
		case message := <-room.broadcast:
			// TODO save the message in the database
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

// IsActive will return iff chatRoom is active
//func (room *ChatRoom) IsActive() bool {
//	return room.isActive
//}
//
//// Activator will set IsActive status of chatRoom to true
//func (room *ChatRoom) Activator() {
//	room.isActive = true
//}

// RegisterClient will register the user to the
func (room *ChatRoom) RegisterClient(client *Client) {
	room.register <- client
}

func CreateNewClient(conn *websocket.Conn, chatRoom *ChatRoom) *Client {
	client := &Client{
		conn: conn,
		room: chatRoom,
		send: make(chan []byte)}
	return client
}
