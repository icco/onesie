package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
		log.Printf("Recieved Message: '%+v'", msgStr)

		if msgStr == "update" {
			// Merge Certs
			files, err := ioutil.ReadDir("/opt/onesie-configs/certs/")
			if err != nil {
				log.Fatalf("Error iterating through files: %+v", err)
			}

			dhparam, err := ioutil.ReadFile(fmt.Sprintf("/opt/onesie-configs/dhparam.pem"))
			if err != nil {
				log.Fatalf("Error reading dhparam: %+v", err)
			}

			for _, file := range files {
				// cat /opt/onesie-configs/certs/$i/{privkey,fullchain}.pem /opt/onesie-configs/dhparam.pem > /opt/onesie-configs/hitch/$i.pem
				log.Printf("Parsing file: %+v", file)
				privkey, err := ioutil.ReadFile(fmt.Sprintf("/opt/onesie-configs/certs/%s/privkey.pem", file.Name()))
				if err != nil {
					log.Fatalf("Error reading privkey: %+v", err)
				}
				fullchain, err := ioutil.ReadFile(fmt.Sprintf("/opt/onesie-configs/certs/%s/fullchain.pem", file.Name()))
				if err != nil {
					log.Fatalf("Error reading fullchain: %+v", err)
				}

				// Write out
				f, err := os.OpenFile(fmt.Sprintf("/opt/onesie-configs/hitch/%s.pem", file.Name()), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					log.Fatalf("Error opening output pem: %+v", err)
				}
				defer f.Close()

				if _, err = f.Write(privkey); err != nil {
					log.Fatalf("Error writing output pem: %+v", err)
				}
				if _, err = f.Write(fullchain); err != nil {
					log.Fatalf("Error writing output pem: %+v", err)
				}
				if _, err = f.Write(dhparam); err != nil {
					log.Fatalf("Error writing output pem: %+v", err)
				}
			}

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
