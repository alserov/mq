package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

const queue = "queue"

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic("failed to dial: " + err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic("failed to init channel: " + err.Error())
	}

	q, err := ch.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		panic("failed to init queue: " + err.Error())
	}

	chDone := make(chan struct{}, 1)

	go func() {
		for i := range 10 {
			b, _ := json.Marshal(i)

			err = ch.PublishWithContext(context.Background(), "", q.Name, false, false, amqp091.Publishing{
				ContentType: "text/plain",
				Body:        b,
			})
			if err != nil {
				fmt.Println("failed to publish: " + err.Error())
			}
		}
	}()

	go func() {
		msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
		if err != nil {
			fmt.Println("failed to consume: " + err.Error())
		}

		for msg := range msgs {
			var i int
			if err = json.Unmarshal(msg.Body, &i); err != nil {
				fmt.Println("consumer error: " + err.Error())
			}

			fmt.Printf("received message: %d\n", i)
			chDone <- struct{}{}
		}
	}()

	for range 10 {
		<-chDone
	}
}
