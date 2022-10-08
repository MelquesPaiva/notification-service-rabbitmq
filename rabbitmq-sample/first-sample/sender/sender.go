package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MelquesPaiva/notification-service-rabbitmq/rabbitmq-sample/first-sample/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	log.Println("Connection succesffully")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declaring a queue that will be sended
	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notification := model.Notification{
		Key:   "grades",
		Value: "9.8",
	}

	byteData, err := json.Marshal(notification)
	failOnError(err, "Failed to convert structure to bytes")

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(byteData),
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %v\n", notification)
}
