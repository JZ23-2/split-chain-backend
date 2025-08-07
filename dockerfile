FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o app

CMD ["./app"]