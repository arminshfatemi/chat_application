package broker

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"notification/database"
	"time"
)

// Message struct to unmarshal sent data from producer
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `bson:"content"`
	SenderID  primitive.ObjectID `bson:"sender_id"`
	RoomID    primitive.ObjectID `bson:"room_id"`
	Timestamp time.Time          `bson:"timestamp"`
}

func RabbitMQConsumer(mongoClient *mongo.Client) {
	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalln("RabbitMQConsumer in dial", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalln("RabbitMQConsumer in defer", err)
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("RabbitMQConsumer in channel", err)
	}

	// declaring the Queue if dont exists
	q, err := ch.QueueDeclare(
		"send_notification",
		false,
		false, // delete
		false,
		false,
		nil,
	)

	// Start Consuming
	messages, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("RabbitMQConsumer in consume", err)
	}

	// here we handle all the received messages form the producer
	for {
		select {
		case msg := <-messages:
			log.Printf("Received a message: %s", msg.Body)
			if err := msg.Ack(false); err != nil {
				log.Fatalln("Ack", err)
			}

			var unmarshalledMessage Message
			err := json.Unmarshal(msg.Body, &unmarshalledMessage)
			if err != nil {
				log.Fatalln("Unmarshal", err)
			}

			// get all the users from the database
			var chatRoom database.Room
			err = mongoClient.Database("chat_app").Collection("rooms").FindOne(
				context.TODO(), bson.M{
					"_id": unmarshalledMessage.RoomID,
				}).Decode(&chatRoom)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					log.Println("RabbitMQ consumer No documents found")
				}
			}

			// get the Client from connection map and send the notification to the channel to send them
			for _, ID := range chatRoom.Members {
				client, ok := database.ConnectedClients[ID.Hex()]
				if !ok {
					log.Println("RabbitMQ consumer no client found", ID.Hex())
					continue
				} else {
					log.Println("RabbitMQ consumer send to channel id of ", ID)
					client.SendTOChan(msg.Body)
				}

			}
		}
	}

}
