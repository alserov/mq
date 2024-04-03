package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

const topic = "topic"

func main() {
	var (
		prod sarama.SyncProducer
		cons sarama.Consumer
		err  error

		chDone = make(chan struct{}, 1)
	)

	// producer
	prodCfg := sarama.NewConfig()
	prodCfg.Producer.Partitioner = sarama.NewRandomPartitioner
	prodCfg.Producer.RequiredAcks = sarama.WaitForAll
	prodCfg.Producer.Return.Successes = true
	prodCfg.ClientID = "1"

	prod, err = sarama.NewSyncProducer([]string{"localhost:29092"}, prodCfg)
	if err != nil {
		panic("failed to init producer: " + err.Error())
	}

	// consumer
	consCfg := sarama.NewConfig()
	consCfg.ClientID = "2"

	cons, err = sarama.NewConsumer([]string{"localhost:29092"}, consCfg)
	if err != nil {
		panic("failed to init consumer: " + err.Error())
	}

	// getting all topic partitions
	partitions, _ := cons.Partitions(topic)

	// consuming partition
	pcons, err := cons.ConsumePartition(topic, partitions[0], sarama.OffsetNewest)
	if nil != err {
		panic("failed to consume: " + err.Error())
	}

	go func() {
		for i := range 10 {
			b, _ := json.Marshal(i)

			if _, _, err = prod.SendMessage(&sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(b),
			}); err != nil {
				fmt.Println("producer error: " + err.Error())
			}
		}
	}()

	go func() {
		for msg := range pcons.Messages() {
			var i int
			if err = json.Unmarshal(msg.Value, &i); err != nil {
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
