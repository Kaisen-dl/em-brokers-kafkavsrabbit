package main

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	name := os.Args[1] 

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	if err := ch.Qos(3, 0, false); err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume("push_queue", "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[push %s] started\n", name)

	for d := range msgs {
		body := string(d.Body)
		fmt.Printf("[push %s] got: %s (redelivered=%v)\n", name, body, d.Redelivered)

		time.Sleep(400 * time.Millisecond)
		d.Ack(false)
	}
}