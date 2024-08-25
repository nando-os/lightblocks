package handler

import (
	"errors"
	"fmt"
	"lightblocks/ordered_map"
	"log"
	"os"
	"sync"
)

const (
	AddItem     = "addItem"
	DeleteItem  = "deleteItem"
	GetItem     = "getItem"
	GetAllItems = "getAllItems"
)

type CommandHandler interface {
	HandleCommandFunc(cmd string, args ...string)
	Close()
}

type handler struct {
	orderedMap ordered_map.Map
	file       *os.File
}

func NewCommandHandler(filename string, orderedMap ordered_map.Map) (CommandHandler, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		msg := fmt.Sprintf("Error opening file: %v", err)
		log.Fatalf(msg)
		return nil, errors.New(msg)

	}

	return &handler{
		orderedMap: orderedMap,
		file:       file,
	}, nil
}

func (s *handler) HandleCommandFunc(cmd string, args ...string) {
	switch cmd {
	case AddItem:
		if len(args) == 2 {
			key, value := args[0], args[1]
			s.orderedMap.Add(key, value)
			// log.Printf("addItem: key=%s, value=%s", key, value)
		} else {
			log.Printf("Invalid arguments: %v", args)
		}
	case DeleteItem:
		if len(args) == 1 || len(args) == 2 {
			key := args[0]
			// log.Printf("deleteItem: key=%s", key)
			s.orderedMap.Delete(key)
		} else {
			log.Printf("Invalid arguments: %v", args)
		}
	case GetItem:
		if len(args) == 1 || len(args) == 2 {
			key := args[0]
			value, found := s.orderedMap.Get(key)
			if found {
				//	log.Printf("getItem: key=%s, value=%s", key, value)
				s.writeStringToFile(fmt.Sprintf("%s %s\n", key, value))
			} else {
				log.Printf("Item not found: key=%s", key)
			}
		} else {
			log.Printf("Invalid arguments: %v", args)
		}
	case GetAllItems:
		items := s.orderedMap.GetAllItems()
		wg := sync.WaitGroup{}
		for _, kv := range items {
			//	log.Printf("getAllItems: key=%s, value=%s", kv.Key, kv.Value)
			wg.Add(1)
			go s.writeToFile(kv, &wg)
		}
		wg.Wait()
	default:
		log.Printf("Unknown command: %s", cmd)
	}

}

// writeToFile writes the output to a file
func (s *handler) writeToFile(kv ordered_map.KeyValue, wg *sync.WaitGroup) {
	defer wg.Done()
	data := fmt.Sprintf("%s %s\n", kv.Key, kv.Value)
	s.writeStringToFile(data)
}

func (s *handler) writeStringToFile(data string) {
	if _, err := s.file.WriteString(data); err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}

func (s *handler) Close() {
	log.Print("closing the handler resources")
	if err := s.file.Close(); err != nil {
		log.Fatalf("Error closing file: %v", err)
	}
}
