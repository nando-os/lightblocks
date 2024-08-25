package rabbit

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"lightblocks/handler"
	"log"
	"sync"
	"time"
)

type Subscriber interface {
	Subscribe(handler handler.CommandHandler) error
	Close()
}

type subscriber struct {
	connectionStr string
	queueName     string
	ctx           context.Context
	wg            *sync.WaitGroup
	connection    *amqp091.Connection
	channel       *amqp091.Channel
}

func NewSubscriber(connStr, queueName string, ctx context.Context) Subscriber {
	return &subscriber{
		connectionStr: connStr,
		queueName:     queueName,
		ctx:           ctx,
		wg:            &sync.WaitGroup{},
	}

}

func (s *subscriber) Subscribe(handler handler.CommandHandler) error {
	conn, err := amqp091.Dial(s.connectionStr)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	s.connection = conn
	log.Print("connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to oepn a channel: %v", err)
	}
	s.channel = ch

	log.Print("opened a channel successfully")

	q, err := ch.QueueDeclare(
		s.queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to regriste a consumer: %v", err)
	}

	go func() {
		sleepTime := 3 * time.Second
		for {
			select {
			case <-msgs:
				for msg := range msgs {
					s.wg.Add(1)
					go func(msg amqp091.Delivery) {
						defer s.wg.Done()
						log.Print("Receiving message ", string(msg.Body))
						var cmd string
						var args1 string
						var args2 string
						fmt.Sscanf(string(msg.Body), "%s %s %s", &cmd, &args1, &args2)
						handler.HandleCommandFunc(cmd, args1, args2)
					}(msg)
				}
			case <-s.ctx.Done():
				log.Print("context done")
				return
			case <-s.channel.NotifyCancel(make(chan string)):
				log.Print("channel cancelled")
				break
			case <-s.channel.NotifyClose(make(chan *amqp091.Error)):
				log.Print("channel closed")
				time.Sleep(sleepTime)
				break
			case <-s.channel.NotifyFlow(make(chan bool)):
				log.Print("channel flow")
				break
			case <-s.channel.NotifyReturn(make(chan amqp091.Return)):
				log.Print("channel return")
				break
			case <-s.channel.NotifyPublish(make(chan amqp091.Confirmation)):
				log.Print("channel publish")
				break
			}
		}
	}()

	log.Print("subscribed to the queue")
	return nil
}

func (s *subscriber) Close() {
	log.Print("closing the subscriber resources")
	// Wait for group to finish before exiting
	s.wg.Wait()
	log.Print("all subscriber goroutines finished")

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
