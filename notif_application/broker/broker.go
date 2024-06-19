package broker

import (
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `bson:"content"`
	SenderID  primitive.ObjectID `bson:"sender_id"`
	RoomID    primitive.ObjectID `bson:"room_id"`
	Timestamp time.Time          `bson:"timestamp"`
}

func RabbitMQConsumer() {
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

		}
	}

}
