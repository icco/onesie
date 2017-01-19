package main

import (
	"context"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

const updateTopicString string = "onesie-updates"
const updateSub string = "onesie-server"

func main() {
	log.Println("Starting")
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
		log.Fatalf("Error getting messages: %+v", err)
	}
	defer it.Stop()

	// Consume 1 messages
	for i := 0; i < 1; i++ {
		msg, err := it.Next()
		if err == iterator.Done {
			log.Println("No more messages.")
			break
		}
		if err != nil {
			log.Printf("Error while getting message: %+v", err)
			break
		}

		msgStr := string(msg.Data)
		log.Printf("Recieved Message: %+v", msgStr)

		if msgStr == "update" {
			// Get hitch PID, send sighup
			out, err := exec.Command("/bin/pidof", "hitch").Output()
			if err != nil {
				log.Printf("Error running pidof: %+v", err)
			}
			for _, pidStr := range strings.Split(string(out), " ") {
				pid, err := strconv.Atoi(strings.TrimSpace(pidStr))
				if err != nil {
					log.Printf("Error parsing string: %+v", err)
					continue
				}
				log.Printf("Sending SIGHUP to %+v", pid)
				syscall.Kill(pid, syscall.SIGHUP)
			}
		}

		msg.Done(true)
	}

	log.Println("Finished.")
}
