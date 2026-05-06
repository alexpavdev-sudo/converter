package queue_conversion

import (
	"context"
	"converter/dto/inner"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"time"
)

const NameQueue = "task_queue"

type RabbitMQConverterQueue struct {
}

func NewRabbitMQConverterQueue() *RabbitMQConverterQueue {
	return &RabbitMQConverterQueue{}
}

func (cq RabbitMQConverterQueue) Push(fileId uint) error {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return errors.New("failed to connect to RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return errors.New("failed to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		NameQueue, // name
		true,      // durability
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	if err != nil {
		return errors.New("failed to declare a queue")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(&inner.MessageDto{FileID: fileId})
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	if err != nil {
		return errors.New("failed to publish a message")
	}

	return nil
}
