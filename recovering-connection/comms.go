package recoveringconnection

import amqp "github.com/rabbitmq/amqp091-go"

type Connection struct {
	conn *amqp.Connection
}
