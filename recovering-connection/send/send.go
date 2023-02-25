package main

import "rabbitmq-docs/recovering-connection/send/comms"

func main() {
	conn := comms.NewConnection("midtrans", "direct", []string{"hello"})
	if err := conn.Connect(); err != nil {
		panic(err)
	}
	if err := conn.BindQueue(); err != nil {
		panic(err)
	}

	m := comms.Message{
		Queue: "hello",
		Body: comms.MessageBody{
			Data: []byte("Hello World"),
		},
	}
	if err := conn.Publish(m); err != nil {
		panic(err)
	}
	// for _, q := range c.queues {
	// 	m := comms.Message{
	// 		Queue: q,
	// 		//set the necessary fields
	// 	}
	// 	if err := conn.Publish(m); err != nil {
	// 		panic(err)
	// 	}
	// }
}

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Panicf("%s: %s", msg, err)
// 	}
// }

// func main() {
// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"hello", // name
// 		true,    // durable
// 		false,   // delete when unused
// 		false,   // exclusive
// 		false,   // no-wait
// 		nil,     // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	body := "Hello World!"
// 	err = ch.PublishWithContext(ctx,
// 		"",     // exchange
// 		q.Name, // routing key
// 		false,  // mandatory
// 		false,  // immediate
// 		amqp.Publishing{
// 			ContentType:  "text/plain",
// 			Body:         []byte(body),
// 			DeliveryMode: 2,
// 		})
// 	failOnError(err, "Failed to publish a message")
// 	log.Printf(" [x] Sent %s\n", body)
// }
