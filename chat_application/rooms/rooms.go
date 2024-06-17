package rooms

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
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
	isActive   bool
	ctx        context.Context
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func CreateNewChatRoom(name string) *ChatRoom {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	return &ChatRoom{
		name:       name,
		isActive:   false,
		ctx:        ctx,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

//func WriteMessage(client *Client) {
//	defer func() {
//		client.room.unregister <- client
//		if err := client.conn.Close(); err != nil {
//			log.Println("WriteMessage: ", err)
//		}
//	}()
//
//	for {
//		select {
//		case message, ok := <-client.send:
//			if !ok {
//				err := client.conn.WriteMessage(websocket.CloseMessage, []byte{})
//				if err != nil {
//					log.Println("WriteMessage inside: ", err)
//				}
//				log.Println("WriteMessage: connection closed")
//				return
//			}
//			client.conn.WriteMessage(websocket.TextMessage, message)
//		}
//	}
//}

//// ReadMessage main function and goroutine to handle sent messages from clients
//func ReadMessage(client *Client) {
//	// if user is disconnected we will remove it from the chatroom
//	defer func() {
//		client.room.unregister <- client
//		if err := client.conn.Close(); err != nil {
//			log.Println("ReadMessage: ", err)
//		}
//	}()
//
//	// TODO: check if need to put a shutdown case in select case
//	// loop that read the message from the user send it to the
//mainLoop:
//	for {
//		select {
//		default:
//			_, message, err := client.conn.ReadMessage()
//			if err != nil {
//				log.Println("ReadMessage: ", err)
//				break mainLoop
//			}
//			client.room.broadcast <- message
//		}
//	}
//}

func WriteMessage(client *Client, message []byte) {
	//defer func() {
	//	client.room.unregister <- client
	//	if err := client.conn.Close(); err != nil {
	//		log.Println("WriteMessage: ", err)
	//	}
	//}()

	err := client.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println(err)
	}

}

// ReadMessage main function and goroutine to handle sent messages from clients
func ReadMessage(client *Client, ctx context.Context) {
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
		case <-ctx.Done():
			if err := client.conn.Close(); err != nil {
				log.Println("in ctx done", err)
			}
			log.Println("gracefully shut down in READ")
			break mainLoop
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

//// Run is function that is responsible for controlling the chatRoom like sending message, user joining and leaving chat
//func (room *ChatRoom) Run() {
//	for {
//		select {
//		// when a user what to join the chat room
//		case client := <-room.register:
//			room.clients[client] = true
//
//		// when a user want to leave the room
//		case client := <-room.unregister:
//			if _, exists := room.clients[client]; exists {
//				delete(room.clients, client)
//				close(client.send)
//			}
//
//		// when a user send a message in the room
//		case message := <-room.broadcast:
//			for client := range room.clients {
//				select {
//				case client.send <- message:
//				default:
//					delete(room.clients, client)
//					close(client.send)
//				}
//			}
//		}
//	}
//}

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

		case message := <-room.broadcast:
			for client := range room.clients {
				go WriteMessage(client, message)
			}

		case <-room.ctx.Done():
			room.isActive = false
			log.Println("gracefully shut down in RUN")
			// TODO complete
			return
		}
	}
}

// ClientsCount will show the count of the clients in the chatRoom
func (room *ChatRoom) ClientsCount() int {
	return len(room.clients)
}

// IsActive will return iff chatRoom is active
func (room *ChatRoom) IsActive() bool {
	return room.isActive
}

// Activator will set IsActive status of chatRoom to true
func (room *ChatRoom) Activator() {
	room.isActive = true
}

// RegisterClient will register the user to the
func (room *ChatRoom) RegisterClient(client *Client) {
	room.register <- client
}

func (room *ChatRoom) ContextGiver() context.Context {
	return room.ctx
}

func CreateNewClient(conn *websocket.Conn, chatRoom *ChatRoom) *Client {
	client := &Client{
		conn: conn,
		room: chatRoom,
		send: make(chan []byte)}
	return client
}

//func ServeWs(room *ChatRoom, c echo.Context) error {
//	conn, err := Upgrader.Upgrade(c.Response(), c.Request(), nil)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "something went wrong")
//	}
//
//	client := &Client{conn: conn, room: room, send: make(chan []byte)}
//	room.register <- client
//	go ReadMessage(client)
//	WriteMessage(client)
//	return nil
//}
