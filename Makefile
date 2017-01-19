all: wrkr status-server

wrkr: wrkr.go
	go build wrkr.go

status-server: status/server.go
	go build -o status-server ./status/server.go

clean:
	rm wrkr
	rm status-server

get-deps:
	go get -u -v cloud.google.com/go/...
	go get -u -v google.golang.org/api/iterator
	go get -u -v github.com/coreos/go-systemd/dbus
	go get -u -v gopkg.in/gin-gonic/gin.v1
