package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
)

const subj = "subj"

func main() {
	chDone := make(chan struct{}, 1)

	conn, err := nats.Connect("localhost:4222")
	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	go func() {
		for i := range 10 {
			b, _ := json.Marshal(i)
			if err = conn.Publish(subj, b); err != nil {
				fmt.Println("failed to publish: " + err.Error())
			}
		}
	}()

	go func() {
		_, err = conn.Subscribe(subj, func(msg *nats.Msg) {
			var i int
			if err = json.Unmarshal(msg.Data, &i); err != nil {
				fmt.Println("failed to unmarshal: " + err.Error())
			}

			fmt.Printf("received a message: %d\n", i)
			chDone <- struct{}{}
		})
	}()

	for range 10 {
		<-chDone
	}
}
