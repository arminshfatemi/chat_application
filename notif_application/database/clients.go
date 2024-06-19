package database

import (
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

var ConnectedClients = make(map[string]*Client)

type Client struct {
	ID               primitive.ObjectID
	conn             *websocket.Conn
	notificationChan chan []byte
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
			client.conn.WriteMessage(websocket.TextMessage, notification)

		}
	}
}
