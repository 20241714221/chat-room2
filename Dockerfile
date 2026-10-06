FROM golang:alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o client ./cmd/client

CMD ["./client"]