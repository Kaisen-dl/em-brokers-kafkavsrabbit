package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func main() {
	instance := os.Args[1] 

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		GroupID:  "notifications-group",
		Topic:    "notifications",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer r.Close()

	fmt.Printf("[consumer %s] started\n", instance)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("[%s] read: %v", instance, err)
		}
		fmt.Printf("[consumer %s] partition=%d key=%s value=%s\n",
			instance, m.Partition, string(m.Key), string(m.Value))
	}
}