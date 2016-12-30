all: wrkr

wrkr: wrkr.go
	go build wrkr.go

clean:
	rm wrkr

get-deps:
	go get -u -v cloud.google.com/go/...
	go get -u -v google.golang.org/api/iterator
