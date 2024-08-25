package main

import (
	"bufio"
	"context"
	"fmt"
	"lightblocks/rabbit"
	"log"
	"os"
	"strings"
	"time"
)

const queueName = "lightblocks"

/*
var (

	hostname = flag.String("host", "localhost", "hostname of the rabbitmq server")
	port     = flag.Int("port", 5672, "port of the rabbitmq server")
	username = flag.String("username", "guest", "username of the rabbitmq server")
	password = flag.String("password", "az00+&z2iop937p", "password of the rabbitmq server")

)

	func init() {
		flag.Parse()
	}
*/
func main() {
	//	// check here number of arguments
	//	if len(os.Args) < 5 {
	//		fmt.Println("Usage: go run main.go <username> <password> <host> <port>")
	//		os.Exit(1)
	//	}

	hostname := "localhost"
	port := 5672
	username := "guest"
	password := "guest"

	// we need username, password, host, port
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, hostname, port)

	publisher, err := rabbit.NewPublisher(connectionString, queueName)
	if err != nil {
		log.Fatalf("Error creating publisher: %v", err)
		return
	}
	defer publisher.Close()

	if len(os.Args) == 2 {
		// If a file is specified, read commands from the file
		filePath := os.Args[1]
		log.Print("Publishing commands from file")
		err := publishCommandsFromFile(publisher, filePath)
		if err != nil {
			log.Fatalf("Error sending commands from file: %v", err)
			return
		}
	} else {
		log.Print("Type in the commands to be sent [cmd + ENTER]:")
		// Otherwise, read commands from standard input
		publishCommandsFromStdin(publisher)
	}
}

func publishCommandsFromFile(publisher rabbit.Publisher, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := scanner.Text()
		if err := publisher.Publish(ctx, command); err != nil {
			log.Fatal(err)
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reding file: %v", err)
	}

	return nil
}

func publishCommandsFromStdin(publisher rabbit.Publisher) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()

		if strings.TrimSpace(command) == "exit" {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := publisher.Publish(ctx, command); err != nil {
			log.Fatal(err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading standard input: %v", err)
	}
}
