// Testing EXCHANGES:
// Exchanges is between the sender and the receiver
// On one side it receives messages from produces and the other side pushes them to queues
// The producer doesn't need to know to which queue the message is being sended

package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/MelquesPaiva/notification-service-rabbitmq/rabbitmq-sample/helper"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	helper.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	log.Println("Connection succesffully")

	ch, err := conn.Channel()
	helper.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		amqp.ExchangeFanout,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)

	helper.FailOnError(err, "Failed to declare an exchange")

	// Declaring a queue that will be sended
	q, err := ch.QueueDeclare(
		"task_queue",
		true, // durable. If rabbitMQ server restart, the queue will survive
		false,
		false,
		false,
		nil,
	)

	helper.FailOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Getting params from terminal
	body := helper.BodyFrom(os.Args)
	err = ch.PublishWithContext(
		ctx,
		"logs", // exchange
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)

	helper.FailOnError(err, "Failed to publish a message")
	log.Printf("[x] Sent %s", body)
}
