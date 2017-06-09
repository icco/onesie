package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
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
		sub, err = pubsubClient.CreateSubscription(context.Background(), updateSub, pubsub.SubscriptionConfig{
			Topic: updateTopic,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	var mu sync.Mutex
	received := 0
	cctx, cancel := context.WithCancel(context.Background())
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		received++
		if received >= 4 {
			cancel()
			msg.Nack()
			return
		}

		msgStr := string(msg.Data)
		fmt.Printf("Got message: %q\n", msgStr)

		if msgStr == "deploy" {
			domain := msg.Attributes["domain"]
			path := msg.Attributes["path"]
			log.Printf("Opening for archive: %s", path)

			// Open google storage client
			client, err := storage.NewClient(ctx)
			if err != nil {
				log.Panicf("Error connecting to Google Storage: %+v", err)
			}
			defer client.Close()
			bkt := client.Bucket("onesie")
			obj := bkt.Object(path)
			r, err := obj.NewReader(ctx)
			if err != nil {
				log.Panicf("Error opening object: %+v", err)
			}
			defer r.Close()

			// Expand into archive
			archive, err := gzip.NewReader(r)
			if err != nil {
				log.Panicf("Error creating gzip reader: %+v", err)
			}
			defer archive.Close()

			// Go through file by file
			tarReader := tar.NewReader(archive)
			buf := make([]byte, 160)
			for {
				header, err := tarReader.Next()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Panicf("Error reading tar: %+v", err)
				}

				path := filepath.Join(domain, header.Name)
				switch header.Typeflag {
				case tar.TypeDir:
					continue
				case tar.TypeReg:
					w := bkt.Object(path).NewWriter(ctx)
					defer w.Close()
					w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
					if filepath.Ext(path) != "" {
						w.ObjectAttrs.ContentType = mime.TypeByExtension(filepath.Ext(path))
					}
					wrtn, err := io.CopyBuffer(w, tarReader, buf)
					if err != nil {
						log.Printf("Error writing data to GCS: %+v", err)
					}
					log.Printf("Wrote %v bytes to %s", wrtn, path)
				default:
					log.Printf("Unable to figure out type: %v (%s)", header.Typeflag, path)
				}
			}
		}

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
					log.Printf("Error reading privkey: %+v", err)
					continue
				}
				fullchain, err := ioutil.ReadFile(fmt.Sprintf("/opt/onesie-configs/certs/%s/fullchain.pem", file.Name()))
				if err != nil {
					log.Printf("Error reading fullchain: %+v", err)
					continue
				}

				// Write out
				f, err := os.OpenFile(fmt.Sprintf("/opt/onesie-configs/hitch/%s.pem", file.Name()), os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
				if err != nil {
					log.Printf("Error opening output pem: %+v", err)
					continue
				}
				defer f.Close()

				if _, err = f.Write(privkey); err != nil {
					log.Printf("Error writing privkey to output pem: %+v", err)
					continue
				}
				if _, err = f.Write(fullchain); err != nil {
					log.Printf("Error writing fullchain to output pem: %+v", err)
					continue
				}
				if _, err = f.Write(dhparam); err != nil {
					log.Printf("Error writing dhparam to output pem: %+v", err)
					continue
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
		msg.Ack()
	})

	if err != nil {
		log.Println(err)
	}

	log.Println("Finished.")
}
