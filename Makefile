GOFILES = $(shell find . -name '*.go')

default: build

workdir:
	mkdir -p workdir

build: workdir/asched

build-native: $(GOFILES)
	go build -o workdir/native-asched .

workdir/asched: $(GOFILES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o workdir/conaschedmd/main.go