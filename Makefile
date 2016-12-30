all: wrkr

wrkr: wrkr.go
	go build wrkr.go

clean:
	rm wrkr
