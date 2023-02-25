package main

import (
	"log"
	"rabbitmq-docs/recovering-connection/receive/comms"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	forever := make(chan bool)
	conn := comms.NewConnection("midtrans", "direct", []string{"hello"})

	if err := conn.Connect(); err != nil {
		panic(err)
	}
	if err := conn.BindQueue(); err != nil {
		panic(err)
	}

	deliveries, err := conn.Consume()
	if err != nil {
		panic(err)
	}
	for q, d := range deliveries {
		go conn.HandleConsumedDeliveries(q, d, messageHandler)
	}
	<-forever
}

func messageHandler(c comms.Connection, q string, deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		m := comms.Message{
			Queue:         q,
			Body:          comms.MessageBody{Data: d.Body, Type: d.Headers["type"]},
			ContentType:   d.ContentType,
			Priority:      d.Priority,
			CorrelationID: d.CorrelationId,
		}
		//handle the custom message
		log.Println("Got message from queue ", m.Queue)
		// d.Ack(false)
	}
}

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Panicf("%s: %s", msg, err)
// 	}
// }

// func DetectError(conn *amqp.Connection) {
// 	// var recvErr *amqp.Error
// 	fmt.Println("DETECT CONNECTION ERROR")
// 	recvErr := make(chan *amqp.Error)
// 	err := <-conn.NotifyClose(recvErr)
// 	if err != nil {
// 		log.Println("koneksi ditutup")
// 		log.Println(err)
// 	}
// 	fmt.Println(<-recvErr)
// }

// func main() {
// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	go DetectError(conn)

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

// 	// q, err := ch.QueueDeclare(
// 	// "task_queue", // name
// 	// true,   // durable
// 	// false,   // delete when unused
// 	// false,   // exclusive
// 	// false,   // no-wait
// 	// nil,     // arguments
// 	// )
// 	failOnError(err, "Failed to declare a queue")

// 	// err = ch.Qos(
// 	// 	1,
// 	// 	0,
// 	// 	false,
// 	// )
// 	// failOnError(err, "Failed to declare a queue")

// 	msgs, err := ch.Consume(
// 		q.Name, // queue
// 		"",     // consumer
// 		false,  // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	)
// 	failOnError(err, "Failed to register a consumer")

// 	var forever chan struct{}

// 	go func() {
// 		for d := range msgs {
// 			// if err:=<-d.; err != nil {

// 			// }

// 			log.Printf("Received a message: %s", d.Body)
// 			dotCount := bytes.Count(d.Body, []byte("."))
// 			t := time.Duration(dotCount)
// 			time.Sleep(t * time.Second)
// 			log.Printf("Done")
// 			// d.Ack(false)
// 		}
// 	}()

// 	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
// 	<-forever

// }
