package comms

import (
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageBody is the struct for the body passed in the AMQP message. The type will be set on the Request header
type MessageBody struct {
	Data []byte
	Type string
}

// Message is the amqp request to publish
type Message struct {
	Queue         string
	ReplyTo       string
	ContentType   string
	CorrelationID string
	Priority      uint8
	Body          MessageBody
}

// Connection is the connection created
type Connection struct {
	name     string
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queues   []string
	err      chan error
}

var (
	connectionPool = make(map[string]*Connection)
)

// NewConnection returns the new connection object
func NewConnection(name, exchange string, queues []string) *Connection {
	if c, ok := connectionPool[name]; ok {
		return c
	}
	c := &Connection{
		name:     name,
		exchange: exchange,
		queues:   queues,
		err:      make(chan error),
	}

	connectionPool[name] = c
	return c
}

// GetConnection returns the connection which was instantiated
func GetConnection(name string) *Connection {
	return connectionPool[name]
}

func (c *Connection) Connect() error {
	var err error
	amqpURI := "amqp://guest:guest@localhost:5672/"
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Error in creating rabbitmq connection with %s : %s", amqpURI, err.Error())
	}
	go func() {
		<-c.conn.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		c.err <- errors.New("Connection Closed")
	}()
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	// fmt.Println("c.exchange in Connect(): ", c.exchange)
	// if err := c.channel.ExchangeDeclare(
	// 	c.name,     // name
	// 	c.exchange, // type
	// 	true,       // durable
	// 	false,      // auto-deleted
	// 	false,      // internal
	// 	false,      // noWait
	// 	nil,        // arguments
	// ); err != nil {
	// 	return fmt.Errorf("Error in Exchange Declare: %s", err)
	// }
	return nil
}

func (c *Connection) BindQueue() error {
	for _, q := range c.queues {
		if _, err := c.channel.QueueDeclare(
			q,     // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		); err != nil {
			return fmt.Errorf("error in declaring the queue %s", err)
		}
		if err := c.channel.QueueBind(q, "midtrans_key_routing", c.exchange, false, nil); err != nil {
			return fmt.Errorf("Queue  Bind error: %s", err)
		}
	}
	return nil
}

// Reconnect reconnects the connection
func (c *Connection) Reconnect() error {
	if err := c.Connect(); err != nil {
		return err
	}
	if err := c.BindQueue(); err != nil {
		return err
	}
	return nil
}

// Publish publishes a request to the amqp queue
func (c *Connection) Publish(m Message) error {
	select { //non blocking channel - if there is no error will go to default where we do nothing
	case err := <-c.err:
		if err != nil {
			c.Reconnect()
		}
	default:
	}

	p := amqp.Publishing{
		Headers:       amqp.Table{"type": m.Body.Type},
		ContentType:   m.ContentType,
		CorrelationId: m.CorrelationID,
		Body:          m.Body.Data,
		ReplyTo:       m.ReplyTo,
		DeliveryMode:  2,
	}

	fmt.Println("Lolos")
	fmt.Println("c.exchange", c.exchange)
	fmt.Println("m.Queue", m.Queue)
	// if err := c.channel.Publish(c.exchange, m.Queue, false, false, p); err != nil {
	if err := c.channel.Publish("", m.Queue, false, false, p); err != nil {
		return fmt.Errorf("Error in Publishing: %s", err)
	}
	return nil
}
