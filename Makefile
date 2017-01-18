all: wrkr server

wrkr: wrkr.go
	go build wrkr.go

server: status/server.go
	go build ./status/server.go

clean:
	rm wrkr
	rm server

get-deps:
	go get -u -v cloud.google.com/go/...
	go get -u -v google.golang.org/api/iterator
	go get -u -v github.com/coreos/go-systemd/dbus
	go get -u -v gopkg.in/gin-gonic/gin.v1
