package lavinmq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client wraps an AMQP connection + channel
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

// New creates a new LavinMQ/AMQP client with retry
func New(url string) *Client {
	c := &Client{url: url}
	c.connect()
	return c
}

func (c *Client) connect() {
	var conn *amqp.Connection
	var err error

	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(c.url)
		if err == nil {
			break
		}
		log.Printf("⚠️  LavinMQ connection attempt %d failed: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ Failed to connect to LavinMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open LavinMQ channel: %v", err)
	}

	c.conn = conn
	c.channel = ch
	log.Println("✅ LavinMQ connected")
}

// DeclareQueue declares a durable queue
func (c *Client) DeclareQueue(name string) (amqp.Queue, error) {
	return c.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
}

// Publish sends a message to a queue
func (c *Client) Publish(ctx context.Context, queue string, body []byte) error {
	_, err := c.DeclareQueue(queue)
	if err != nil {
		return err
	}

	return c.channel.PublishWithContext(ctx,
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

// Consume registers a consumer for a queue
func (c *Client) Consume(queue string, consumer string) (<-chan amqp.Delivery, error) {
	_, err := c.DeclareQueue(queue)
	if err != nil {
		return nil, err
	}

	return c.channel.Consume(
		queue,    // queue
		consumer, // consumer tag
		false,    // auto-ack (manual ack for reliability)
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
}

// Close gracefully closes the AMQP connection
func (c *Client) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	log.Println("✅ LavinMQ connection closed")
}
