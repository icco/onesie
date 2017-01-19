package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/coreos/go-systemd/dbus"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	r := gin.Default()
	r.GET("/.well-known/onesie/status.json", func(c *gin.Context) {
		conn, err := dbus.NewSystemdConnection()
		if err != nil {
			log.Fatalf("Error openning connection: %+v", err)
		}
		defer conn.Close()

		units, err := conn.ListUnits()
		if err != nil {
			log.Fatalf("Error listing units: %+v", err)
		}

		services := gin.H{}
		for _, s := range units {
			if strings.Contains(s.Name, ".service") {
				services[s.Name] = s.ActiveState
			}
		}

		c.JSON(200, gin.H{
			"status":   "ok",
			"services": services,
		})
	})

	port := "9090"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	r.Run(fmt.Sprintf(":%s", port))
}
