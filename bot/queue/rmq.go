package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/url"
	"recipebot/environment"
	"time"
)

type RMQueue struct {
	rmqConn *amqp091.Connection
	channel *amqp091.Channel
	q       *amqp091.Queue
}

func (queue *RMQueue) Configure(config *environment.RmqConfig) error {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.User,
		url.QueryEscape(config.Password),
		config.Host,
		config.Port,
		config.VirtualHost)

	conn, err := amqp091.Dial(connectionString)
	if err != nil {
		return err
	}
	queue.rmqConn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	queue.channel = ch

	q, err := queue.channel.QueueDeclare(
		config.OutputQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	queue.q = &q
	return nil
}

func (queue *RMQueue) SendMessage(message interface{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = queue.channel.PublishWithContext(ctx,
		"",
		queue.q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        jsonMsg,
		})
	if err != nil {
		return err
	}
	log.Printf(" [x] Sent %s %s\n", queue.q.Name, jsonMsg)
	return err
}

func (queue *RMQueue) Close() error {
	log.Println("Disconnecting RabbitMQ")
	qErr := queue.channel.Close()
	connErr := queue.rmqConn.Close()
	return errors.Join(qErr, connErr)
}
