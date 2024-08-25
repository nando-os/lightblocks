package rabbit

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type Publisher interface {
	Publish(ctx context.Context, message string) error
	Close()
}

type publisher struct {
	queueName  string
	connection *amqp091.Connection
	channel    *amqp091.Channel
	queue      *amqp091.Queue
}

func NewPublisher(connectionStr, queueName string) (Publisher, error) {
	log.Print("connecting to RabbitMQ")
	conn, err := amqp091.Dial(connectionStr)
	if err != nil {
		err := fmt.Errorf("failed to connect to RabbitMQ: %v", err)
		log.Fatalf(err.Error())
		return nil, err
	}
	log.Print("connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		err := fmt.Errorf("failed to open a channel: %v", err)
		log.Fatalf(err.Error())
		return nil, err
	}
	log.Print("channel successfully opened")

	queue, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	return &publisher{
		connection: conn,
		channel:    ch,
		queue:      &queue,
	}, nil
}

func (p *publisher) Publish(ctx context.Context, command string) error {
	if p.channel == nil || p.channel.IsClosed() {
		return fmt.Errorf("channel is not open")
	}

	if err := p.channel.PublishWithContext(ctx,
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(command),
		},
	); err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}
	log.Printf("Sent command: %s", command)
	return nil
}

func (s *publisher) Close() {
	log.Print("closing the publisher resources")

	if !s.channel.IsClosed() {
		err := s.channel.Close()
		if err != nil {
			log.Fatalf("Failed to close the channel: %v", err)
		}
		log.Print("channel closed")
	}

	if s.connection != nil && !s.connection.IsClosed() {
		err := s.connection.Close()
		if err != nil {
			log.Fatalf("Failed to close the connection: %v", err)
		}
		log.Print("connection closed")
		return
	}
	log.Print("found no open connection to close")

}
