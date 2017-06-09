all: wrkr status-server

wrkr: wrkr.go
	go build wrkr.go

status-server: status/server.go
	go build -o status-server ./status/server.go

clean:
	rm wrkr
	rm status-server

get-deps:
	go get -u github.com/tools/godep
	go get -u cloud.google.com/go/...
	go get -u github.com/coreos/go-systemd/dbus
	go get -u golang.org/x/net/context
	go get -u gopkg.in/gin-gonic/gin.v1
	${GOPATH}/bin/godep save
