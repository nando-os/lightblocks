package main

import (
	"context"
	"fmt"
	"lightblocks/handler"
	"lightblocks/ordered_map"
	"lightblocks/rabbit"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	rabbitHostname := getEnvironmentVariableOrDefault("RABBITMQ_HOST", "localhost")
	rabbitPort := getEnvironmentVariableOrDefaultInt("RABBITMQ_PORT", 5672)
	rabbitUsername := getEnvironmentVariableOrDefault("RABBITMQ_USER", "guest")
	rabbitPassword := getEnvironmentVariableOrDefault("RABBITMQ_PASS", "guest")
	rabbitQueueName := getEnvironmentVariableOrDefault("RABBITMQ_QUEUE", "lightblocks")
	filename := "OUTPUT/" + getEnvironmentVariableOrDefault("OUTPUT_FILENAME", "server-output.txt")

	// context can be passed to notify cancellation
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	handler, err := handler.NewCommandHandler(filename, ordered_map.NewConcurrentOrderedMap())
	if err != nil {
		log.Fatalf("Error creating handler: %v", err)
		return
	}

	// Subscribe to queue
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%d/", rabbitUsername, rabbitPassword, rabbitHostname, rabbitPort)
	subscriber := rabbit.NewSubscriber(connectionString, rabbitQueueName, ctx)

	server := NewServer(handler, subscriber)
	server.Start(ctx)
	server.WaitForShutdown(ctx)
}

type Server struct {
	handler    handler.CommandHandler
	subscriber rabbit.Subscriber
}

func NewServer(handler handler.CommandHandler, subscriber rabbit.Subscriber) *Server {
	return &Server{handler: handler, subscriber: subscriber}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.subscriber.Subscribe(s.handler); err != nil {
		err := fmt.Errorf("failed to subscribe: %v", err)
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *Server) WaitForShutdown(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Print("server start shutdown...")
			s.subscriber.Close()
			log.Print("server exited")
			return
		}
	}
}

func getEnvironmentVariableOrDefault(envVar, defaultVar string) string {
	if value, ok := os.LookupEnv(envVar); ok {
		return value
	}
	log.Printf("Environment variable %s not set, using default value", envVar)
	return defaultVar
}

func getEnvironmentVariableOrDefaultInt(envVar string, defaultVar int) int {
	if value, ok := os.LookupEnv(envVar); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	log.Printf("Environment variable %s not set or invalid, using default value", envVar)
	return defaultVar
}
