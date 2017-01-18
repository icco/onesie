all: wrkr

wrkr: wrkr.go
	go build wrkr.go

server: status/server.go
	go build ./status/server.go

clean:
	rm wrkr

get-deps:
	go get -u -v cloud.google.com/go/...
	go get -u -v google.golang.org/api/iterator
	go get -u -v github.com/coreos/go-systemd/dbus
