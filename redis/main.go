package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const channel = "channel"

func main() {
	conn := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer conn.Close()

	if err := conn.Ping(context.Background()).Err(); err != nil {
		panic("failed to ping: " + err.Error())
	}

	chDone := make(chan struct{}, 1)
	sub := conn.Subscribe(context.Background(), channel)

	go func() {
		time.Sleep(time.Millisecond * 50)
		for range 10 {
			conn.Publish(context.Background(), channel, "message")
		}
	}()

	go func() {
		for {
			msg, err := sub.ReceiveMessage(context.Background())
			if err != nil {
				fmt.Println("failed to receive message")
			}

			fmt.Println("received message: " + msg.Payload)
			chDone <- struct{}{}
		}
	}()

	for range 10 {
		<-chDone
	}
}
