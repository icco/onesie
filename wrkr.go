package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

const updateTopicString string = "onesie-updates"
const updateSub string = "onesie-server"

func main() {
	log.Println("Hello")
	pubsubClient, err := pubsub.NewClient(context.Background(), "940380154622")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	updateTopic := pubsubClient.Topic(updateTopicString)
	sub := pubsubClient.Subscription(updateSub)
	b, err := sub.Exists(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if !b {
		sub, err = pubsubClient.CreateSubscription(context.Background(), updateSub, updateTopic, 0, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Construct the iterator
	it, err := sub.Pull(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer it.Stop()

	// Consume 1 messages
	for i := 0; i < 1; i++ {
		msg, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
			break
		}

		log.Print("got message: ", string(msg.Data))
		msg.Done(true)
	}
}