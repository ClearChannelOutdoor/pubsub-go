package main

import (
	"cloud.google.com/go/pubsub"
	"fmt"
	pubsub_go "github.com/clearchanneloutdoor/pubsub-go"
	"sync"
)

const (
	projectID    = "develop-applications-75118"
	settingsFile = "./pubsub-dev.json"
	subscription = "pubsub-go-module-test-sub"
	topic        = "pubsub-go-module-test"
)

func main() {
	config := pubsub_go.Config{
		ProjectID:              projectID,
		IsLocal:                true,
		ServiceAccountFilePath: settingsFile,
	}

	ps, err := pubsub_go.NewPubSub(config)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go receive(ps, &wg)

	message := pubsub_go.Message{
		Attributes: nil,
		Message:    "hello world",
		Topic:      topic,
	}

	if err := ps.Publish(message); err != nil {
		panic(err)
	}

	fmt.Println("Published message")

	wg.Wait()

	fmt.Println("Done")
}

func receive(ps *pubsub_go.PubSub, wg *sync.WaitGroup) {
	messages := make(chan *pubsub.Message)
	go func() {
		for {
			msg := <-messages
			fmt.Printf("Got message: %s\n", string(msg.Data))
			msg.Ack()
			wg.Done()
		}
	}()

	if err := ps.Receive(subscription, messages); err != nil {
		panic(err)
	}

	fmt.Println("Received message")
}
