package main

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	psb "github.com/clearchanneloutdoor/pubsub-go/pkg"
)

const (
	projectID    = "pubsub-go-module-test-project"
	subscription = "pubsub-go-module-test-sub"
	topic        = "pubsub-go-module-test-topic"
)

func main() {
	// create a new PubSub client
	client, err := psb.NewPubSub(context.Background(), psb.Options(projectID))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			panic(err)
		}
	}()

	// create a subscription to the topic
	if err := client.CreateSubscription(topic, subscription, ""); err != nil {
		panic(err)
	}

	// create a waitgroup to wait for the message to be received
	var wg sync.WaitGroup
	wg.Add(1)

	// publish an example message
	fmt.Printf("publishing message...\n")
	go receive(client, &wg)
	if err := client.Publish(topic, "hello world"); err != nil {
		panic(err)
	}

	// wait for the message to be received
	wg.Wait()
	fmt.Printf("publish and receive completed")
}

func receive(ps *psb.PubSub, wg *sync.WaitGroup) {
	messages := make(chan *pubsub.Message)
	go func() {
		for {
			msg := <-messages
			fmt.Printf("received message: %s\n", string(msg.Data))

			msg.Ack()
			wg.Done()
		}
	}()

	if err := ps.Receive(subscription, messages); err != nil {
		panic(err)
	}
}
