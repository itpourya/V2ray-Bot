FROM golang:1.19.0-alpine3.16 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o bin .

FROM alpine:3.13.1

CMD ["/app/bin"]
