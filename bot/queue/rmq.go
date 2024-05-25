package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/url"
	"recipebot/config"
	"time"
)

const (
	Name = "RabbitMQ"
)

type RMQueue struct {
	rmqConn *amqp091.Connection
	channel *amqp091.Channel
	q       *amqp091.Queue
	configs *config.Config
}

func (queue *RMQueue) Configure(configs config.Config) {
	queue.configs = &configs
	user := configs.GetString(config.RABBITMQ_USER)
	password := url.QueryEscape(configs.GetString(config.RABBITMQ_PASSWORD))
	host := configs.GetString(config.RABBITMQ_HOST)
	port := configs.GetString(config.RABBITMQ_PORT)
	virtHost := configs.GetString(config.RABBITMQ_VIRT_HOST)

	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, virtHost)

	conn, err := amqp091.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")
	queue.rmqConn = conn

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	queue.channel = ch

	q, err := queue.channel.QueueDeclare(
		config.RABBITMQ_QUEUE_NAME, // name
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	queue.q = &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (queue *RMQueue) SendMessage(message interface{}) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonMsg, err := json.Marshal(message)
	failOnError(err, "Unable to serialize message")

	err = queue.channel.PublishWithContext(ctx,
		"",
		queue.q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        jsonMsg,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", jsonMsg)
}

func (queue *RMQueue) Close() error {
	log.Println("Disconnecting RabbitMQ")
	qErr := queue.channel.Close()
	connErr := queue.rmqConn.Close()
	return errors.Join(qErr, connErr)
}
