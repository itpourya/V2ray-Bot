FROM golang:1.19.0-alpine3.16 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o /app/app .

FROM alpine:3.13.1
RUN apk add --no-cache tzdata
ENV TZ=Asia/Tehran
WORKDIR /app
COPY --from=builder /app/app .

CMD ["/app/app"]
