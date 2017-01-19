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
		conn, err := dbus.NewSystemdConnection()
		if err != nil {
			log.Fatalf("Error openning connection: %+v", err)
		}
		defer conn.Close()

		units, err := conn.ListUnits()
		if err != nil {
			log.Fatalf("Error listing units: %+v", err)
		}
		log.Printf("Units: %+v", units)

		// units: [{Name:systemd-vconsole-setup.service Description:systemd-vconsole-setup.service LoadState:not-found ActiveState:inactive SubState:dead Followed: Path:/org/freedesktop/systemd1/unit/systemd_2dvconsole_2dsetup_2eservice JobId:0 JobType: JobPath:/}]
		services := gin.H{}
		for _, s := range units {
			services[s.Name] = s.ActiveState
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
