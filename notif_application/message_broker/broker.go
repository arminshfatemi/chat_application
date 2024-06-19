package message_broker

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
)

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
		"hello",
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

	// here we handle all the received messages
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
