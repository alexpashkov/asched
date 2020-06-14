GOFILES = $(shell find . -name '*.go')

default: build

bin:
	mkdir -p bin

build: bin/asched

build-native: $(GOFILES)
	go build -o bin/native-asched .

bin/asched: $(GOFILES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/asched cmd/main.go