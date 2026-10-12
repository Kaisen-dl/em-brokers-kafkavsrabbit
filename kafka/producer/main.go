package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "notifications",
		Balancer: &kafka.Hash{}, 
	}
	defer w.Close()

	userIDs := []string{"user-1", "user-2", "user-3", "user-4"}

	for i := 0; i < 3; i++ {
		for _, uid := range userIDs {
			msg := kafka.Message{
				Key:   []byte(uid),
				Value: []byte(fmt.Sprintf("notification.send for %s #%d", uid, i)),
			}
			if err := w.WriteMessages(context.Background(), msg); err != nil {
				log.Fatalf("write: %v", err)
			}
			fmt.Printf("sent key=%s value=%s\n", uid, string(msg.Value))
			time.Sleep(300 * time.Millisecond)
		}
	}
}