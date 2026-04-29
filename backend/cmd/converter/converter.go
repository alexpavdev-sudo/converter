package main

import (
	"converter/app"
	"converter/dto/inner"
	"converter/services/converter"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func exit(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	app.Init(true)
	defer app.App().DeInit()

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	exit(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	exit(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue",
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	exit(err, "Failed to declare a queue")

	err = ch.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	exit(err, "Failed to set QoS")

	deliveries, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	exit(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range deliveries {
			go func(delivery amqp.Delivery) {
				var msg inner.MessageDto
				err := json.Unmarshal(delivery.Body, &msg)
				if err != nil {
					log.Printf("error: %s", err)
					return
				}

				err = converter.NewConverter(msg.FileID).Run()
				if err != nil {
					log.Printf(err.Error())
				}

				delivery.Ack(false)
			}(d)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
