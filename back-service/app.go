package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"log"
	"net/url"
)

const (
	RABBITMQ_USER      = "RABBITMQ_USER"
	RABBITMQ_PASSWORD  = "RABBITMQ_PASSWORD"
	RABBITMQ_HOST      = "RABBITMQ_HOST"
	RABBITMQ_PORT      = "RABBITMQ_PORT"
	RABBITMQ_VIRT_HOST = "RABBITMQ_VIRT_HOST"
)

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	user := viper.GetString(RABBITMQ_USER)
	password := url.QueryEscape(viper.GetString(RABBITMQ_PASSWORD))
	host := viper.GetString(RABBITMQ_HOST)
	port := viper.GetString(RABBITMQ_PORT)
	virtHost := viper.GetString(RABBITMQ_VIRT_HOST)

	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, virtHost)

	conn, err := amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
