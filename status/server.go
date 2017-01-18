package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-systemd/dbus"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	r := gin.Default()
	r.GET("/.well-known/onesie/status.json", func(c *gin.Context) {
		conn, err := dbus.New()
		if err != nil {
			log.Fatalf("Error openning connection: %+v", err)
		}
		defer conn.Close()

		units, err := conn.ListUnits()
		if err != nil {
			log.Fatalf("Error listing units: %+v", err)
		}
		log.Printf("Units: %+v", units)

		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	port := "9090"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	r.Run(fmt.Sprintf(":%s", port))
}
