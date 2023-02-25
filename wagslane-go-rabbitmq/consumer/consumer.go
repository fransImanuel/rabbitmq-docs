package main

import (
	"log"

	"github.com/wagslane/go-rabbitmq"
)

func main() {
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed: %v", string(d.Body))
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		"hello", //queue name here
		// rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
		// rabbitmq.WithConsumerOptionsQueueArgs( rabbitmq.Table{} ),
		// rabbitmq.WithConsumerOptionsExchangeName("events"),
		// rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsQueueDurable,
	)

	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()
}
