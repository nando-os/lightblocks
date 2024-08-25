package handler

import (
	"github.com/stretchr/testify/assert"
	"lightblocks/ordered_map"
	"log"
	"os"
	"testing"
)

func TestHandleCommandFunc(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up after test

	mockMap := ordered_map.NewConcurrentOrderedMap()

	// Create a new handler
	handler, err := NewCommandHandler(tmpfile.Name(), mockMap)
	if err != nil {
		t.Fatalf("Error creating handler: %v", err)
	}

	// Define test cases
	tests := []struct {
		cmd  string
		args []string
	}{
		{AddItem, []string{"key1", "value1"}},
		{DeleteItem, []string{"key1"}},
		{GetItem, []string{"key1"}},
		{GetAllItems, []string{}},
		{"unknownCommand", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			handler.HandleCommandFunc(tt.cmd, tt.args...)

			switch tt.cmd {
			case AddItem:
				value, found := mockMap.Get("key1")
				assert.True(t, found)
				assert.Equal(t, "value1", value)
			case DeleteItem:
				_, found := mockMap.Get("key1")
				assert.False(t, found)
			case GetItem:
				// Check the log output or file content if necessary

			case GetAllItems:
				// Check the log output or file content if necessary
			case "unknownCommand":
				// Check the log output for unknown command
			}
		})
	}
}
