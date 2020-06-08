FROM golang:1.14.4

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY graph graph
COPY internal internal
COPY cmd cmd

RUN go build -v -o asched cmd/main.go

CMD ["./asched"]
