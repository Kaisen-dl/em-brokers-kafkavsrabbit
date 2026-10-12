package main

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
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

	if err := ch.ExchangeDeclare("notifications", "topic", true, false, false, false, nil); err != nil {
		log.Fatal(err)
	}

	if err := ch.ExchangeDeclare("notifications.dlx", "direct", true, false, false, false, nil); err != nil {
		log.Fatal(err)
	}

	_, err = ch.QueueDeclare("email.retry", true, false, false, false, amqp.Table{
		"x-message-ttl":             int32(5000),
		"x-dead-letter-exchange":    "notifications",
		"x-dead-letter-routing-key": "notification.email",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = ch.QueueDeclare("push.retry", true, false, false, false, amqp.Table{
		"x-message-ttl":             int32(5000),
		"x-dead-letter-exchange":    "notifications",
		"x-dead-letter-routing-key": "notification.push",
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := ch.QueueBind("email.retry", "email.retry", "notifications.dlx", false, nil); err != nil {
		log.Fatal(err)
	}
	if err := ch.QueueBind("push.retry", "push.retry", "notifications.dlx", false, nil); err != nil {
		log.Fatal(err)
	}

	_, err = ch.QueueDeclare("email_queue", true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    "notifications.dlx",
		"x-dead-letter-routing-key": "email.retry",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = ch.QueueDeclare("push_queue", true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    "notifications.dlx",
		"x-dead-letter-routing-key": "push.retry",
	})
	if err != nil {
		log.Fatal(err)
	}

	ch.QueueBind("email_queue", "notification.email", "notifications", false, nil)
	ch.QueueBind("push_queue", "notification.push", "notifications", false, nil)

	ctx := context.Background()
	for i := 1; i <= 10; i++ {
		key := "notification.email"
		if i%2 == 0 {
			key = "notification.push"
		}
		body := fmt.Sprintf("msg #%d key=%s", i, key)
		err = ch.PublishWithContext(ctx, "notifications", key, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[publisher] sent %s\n", body)
		time.Sleep(500 * time.Millisecond)
	}
}