package database

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

var ConnectedClients = make(map[string]*Client)

type Room struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	Name      string               `bson:"name"`
	Members   []primitive.ObjectID `bson:"members"`
	CreatedAt time.Time            `bson:"created_at"`
}

type Client struct {
	ID               primitive.ObjectID
	conn             *websocket.Conn
	notificationChan chan Message
}

func CreateNewClient(conn *websocket.Conn, id primitive.ObjectID) *Client {
	return &Client{
		ID:               id,
		conn:             conn,
		notificationChan: make(chan Message),
	}
}

func (client *Client) SendTOChan(msg Message) {
	client.notificationChan <- msg
}

func (client *Client) Run(mongoClient *mongo.Client) {
	defer func() {
		if err := client.conn.Close(); err != nil {
			log.Printf("error closing connection to client: %v", err)
		}
		delete(ConnectedClients, client.ID.Hex())
	}()

	for {
		select {
		case notification := <-client.notificationChan:
			notif := CreateNewNotificationsMessage(&notification, client.ID)
			if err := notif.InsertNotificationInDatabase(mongoClient); err != nil {
				log.Fatalln("error inserting notification: ", err)
			}

			marshaledNotif, err := json.Marshal(&notif)
			if err != nil {
				log.Fatalln("error marshaling notification: ", err)
			}

			if err := client.conn.WriteMessage(websocket.TextMessage, marshaledNotif); err != nil {
				log.Printf("error sending notification to client: %v", err)
			}
		}
	}
}
