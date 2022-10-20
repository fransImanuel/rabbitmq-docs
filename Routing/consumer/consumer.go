package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	con,err := amqp.Dial("amqp:localhost")
	failOnError(err, "Failed to Dial")
	defer con.Close()

	ch,err := con.Channel()
	failOnError(err, "Failed to create channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct2", //name
		"direct", //kind
		false, //durable
		false,//autoDelete
		false, //internal
		false, //nowait
		nil,//args
	)
	failOnError(err, "Failed to exchange declare")

	q,err := ch.QueueDeclare(
		"", //name
		true, //durable
		false, //autodelete
		false, //exclusive
		false, //nowait
		nil, //args
	)
	failOnError(err, "Failed to create queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
				q.Name, "logs_direct2", s)
		err = ch.QueueBind(
				q.Name,        // queue name
				s,             // routing key
				"logs_direct2", // exchange
				false,
				nil)
		failOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name, //queue
		"",//consumer
		false,//auto ack
		false,//exclusive
		false,//nolocal
		false,//nowait
		nil, //args
	)
	failOnError(err, "Failed to consume")

	var forever chan struct{}

	go func() {
		for d:=range msgs{
			log.Printf("[x] %s", d.Body)
			// d.Ack(false)
		}
	}()

	log.Printf(" [x] Waiting for logs. To exit press ctrl+c")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}