package broker

import (
	"chatRoom/models"
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

func NotificationProducer(notificationProducerChan chan models.Message) {
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

	// here we handle all the received messages form the chats and send them to the Queue
	for {
		select {
		case message := <-notificationProducerChan:

			body, err := json.Marshal(message)
			if err != nil {
				log.Fatalln("RabbitMQConsumer in json marshal", err)
			}
			err = ch.PublishWithContext(
				context.Background(),
				"",
				q.Name,
				false,
				false,
				amqp091.Publishing{
					ContentType: "application/json",
					Body:        body,
				})
			if err != nil {
				log.Fatalln("RabbitMQConsumer in publish", err)
			}
		}
	}
}
