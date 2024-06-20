package broker

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"notification/database"
)

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
			// unMarshal the msg to use the RoomID to get the members
			var unmarshalledMessage database.Message
			err := json.Unmarshal(msg.Body, &unmarshalledMessage)
			if err != nil {
				log.Fatalln("Unmarshal", err)
			}

			// because we need members of the room to send them notification we get the room from database
			// to send them notifications, NOTE: we need their ObjectID
			var chatRoom database.Room
			err = mongoClient.Database("chat_app").Collection("rooms").FindOne(
				context.TODO(), bson.M{
					"_id": unmarshalledMessage.RoomID,
				}).Decode(&chatRoom)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					log.Fatalln("RabbitMQ consumer No documents found")
				}
			}

			// iterate throw the members of the chatRoom to use their channel and sending the notification to the
			// Run goroutine and save the notification in database
			for _, ID := range chatRoom.Members {
				client, ok := database.ConnectedClients[ID.Hex()]
				if !ok {
					log.Println("RabbitMQ consumer no client found", ID.Hex())
					continue
				}
				client.SendTOChan(unmarshalledMessage)
			}
			// acknowledgement of the message
			log.Printf("Received a message: %s", unmarshalledMessage)
			if err := msg.Ack(false); err != nil {
				log.Fatalln("Ack", err)
			}
		}
	}

}
