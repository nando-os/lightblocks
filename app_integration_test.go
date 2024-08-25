package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	handler2 "lightblocks/handler"
	"lightblocks/ordered_map"
	"lightblocks/rabbit"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestCommandHandler_HandleCommandFunc(t *testing.T) {
	rabbitHostname := getEnvironmentVariableOrDefault("RABBITMQ_HOST", "localhost")
	rabbitPort := getEnvironmentVariableOrDefaultInt("RABBITMQ_PORT", 5672)
	rabbitUsername := getEnvironmentVariableOrDefault("RABBITMQ_USER", "guest")
	rabbitPassword := getEnvironmentVariableOrDefault("RABBITMQ_PASS", "guest")
	rabbitQueueName := getEnvironmentVariableOrDefault("RABBITMQ_QUEUE", "lightblocks")
	filename := "OUTPUT/" + getEnvironmentVariableOrDefault("OUTPUT_FILENAME", "it-test-output.txt")

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	tests := []struct {
		name          string
		command       string
		key           string
		value         string
		addedToMap    bool
		sizeOfMap     int
		fileEntrySize int
	}{
		{
			name:          "test addItem command 1",
			command:       "addItem",
			key:           "Johan",
			value:         "Cruijf",
			addedToMap:    true,
			sizeOfMap:     1,
			fileEntrySize: 0,
		}, {
			name:          "test addItem command 2",
			command:       "addItem",
			key:           "Dennis",
			value:         "Bergkamp",
			addedToMap:    true,
			sizeOfMap:     2,
			fileEntrySize: 0,
		}, {
			name:          "test getItem",
			command:       "getItem",
			key:           "Dennis",
			value:         "",
			sizeOfMap:     2,
			fileEntrySize: 1,
		}, {
			name:          "test deleteItem",
			command:       "deleteItem",
			key:           "Dennis",
			value:         "",
			sizeOfMap:     1,
			fileEntrySize: 1,
		}, {
			name:          "test addItem command 3",
			command:       "addItem",
			key:           "Patrick",
			value:         "Kluivert",
			addedToMap:    true,
			sizeOfMap:     2,
			fileEntrySize: 1,
		},
		{
			name:          "test getAllItems",
			command:       "getAllItems",
			key:           "",
			value:         "",
			sizeOfMap:     2,
			fileEntrySize: 3,
		},
	}

	// start the rabbitmq subscriber
	orderedMap := ordered_map.NewConcurrentOrderedMap()
	handler, err := handler2.NewCommandHandler(filename, orderedMap)
	if err != nil {
		t.Fatalf("Error creating handler: %v", err)
	}
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%d/", rabbitUsername, rabbitPassword, rabbitHostname, rabbitPort)
	subscriber := rabbit.NewSubscriber(connectionString, rabbitQueueName, ctx)

	if err := subscriber.Subscribe(handler); err != nil {
		t.Fatal(err)
	}
	publisher, err := rabbit.NewPublisher(connectionString, rabbitQueueName)
	if err != nil {
		t.Fatalf("failed to create publisher: %v", err)
	}
	defer publisher.Close()

	publisher.Publish(ctx, "clear")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			command := fmt.Sprintf("%s %s %s", tt.command, tt.key, tt.value)
			if err = publisher.Publish(ctx, command); err != nil {
				t.Fatalf("failed to publish message: %v", err)
			}
			time.Sleep(1 * time.Second)

			assert.Equal(t, tt.sizeOfMap, orderedMap.Size())

			if tt.command == "addItem" {
				value, bool := orderedMap.Get(tt.key)
				assert.True(t, bool)
				assert.Equal(t, tt.value, value)
			}

			if tt.command == "deleteItem" {
				_, bool := orderedMap.Get(tt.key)
				assert.False(t, bool)
			}

			if tt.fileEntrySize > 0 {
				// read the contents of the file into a slice of strings
				fileContents := ReadFileContents(filename)
				assert.Equal(t, tt.fileEntrySize, len(fileContents))
			}

		})
	}

	// assert file content
	fileContents := ReadFileContents(filename)
	assert.Equal(t, 3, len(fileContents))
	assert.Equal(t, "Dennis Bergkamp", fileContents[0])
	assert.Equal(t, "Johan Cruijf", fileContents[1])
	assert.Equal(t, "Patrick Kluivert", fileContents[2])

}

func ReadFileContents(s string) []string {
	file, err := os.Open(s)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines

}
