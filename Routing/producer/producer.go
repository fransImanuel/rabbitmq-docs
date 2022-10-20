package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	con,err :=amqp.Dial("amqp:localhost")
	failOnError(err, "Failed to Dial")
	defer con.Close()

	ch,err := con.Channel()
	failOnError(err, "Fail to create channel")
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
	failOnError(err, "Fail to declare exchange")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()


	body := bodyFrom(os.Args)
	// body := []byte{"Hello World"}
	err = ch.PublishWithContext(ctx,
		"logs_direct2", //name
		severityFrom(os.Args), //key
		false, //durable
		false,//mandatory
		amqp.Publishing{
			ContentType: "plain/text",
			Body: []byte(body),
		})
	failOnError(err, "Failed to publish msg")
	
	log.Printf(" [x] Sent %s", body)
}


func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
			s = "hello"
	} else {
			s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
			s = "info"
	} else {
			s = os.Args[1]
	}
	return s
}