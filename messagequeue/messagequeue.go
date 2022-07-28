package messagequeue

import (
	"bidoof/helpers"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/streadway/amqp"
)

const RETRY_COUNT = 3

var mqUrl = os.Getenv("MESSAGE_QUEUE_URL")
var mqQueueName = os.Getenv("MESSAGE_QUEUE_NAME")

func Connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial(mqUrl)
	retryCount := 0
	//When trying to run this with other containers, I noticed this would try to start
	//before the others were fully functional, so a few retries might do the trick here

	//Probably 100 better ways of doing this...
	for err != nil && retryCount < RETRY_COUNT {
		retryCount++
		helpers.RetrySleep("RabbitMQ")
		conn, err = amqp.Dial(mqUrl)
	}

	return conn, err
}

func Channel(conn *amqp.Connection) (*amqp.Channel, error) {
	//this should be in its own retry-loop function
	return conn.Channel()
}

func Consumer(channel *amqp.Channel) (<-chan amqp.Delivery, error) {

	return channel.Consume(
		mqQueueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
}
