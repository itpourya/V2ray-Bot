# Stage 1: Build the Go application
FROM golang:1.20 AS builder

WORKDIR /app

# Copy go.mod and go.sum files for dependency resolution
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application as a static binary
RUN CGO_ENABLED=0 go build -o telebot . && ls -l /app

# Stage 2: Create a minimal image for the application
FROM alpine:latest

# Install timezone data
RUN apk add --no-cache tzdata

# Set the timezone environment variable
ENV TZ=Asia/Tehran

# Set the working directory
WORKDIR /root/

# Copy the .env file and set proper permissions
COPY .env .
RUN chmod 644 .env

# Copy the built application from the builder stage
COPY --from=builder /app/telebot .

# Ensure the binary is executable
RUN chmod +x ./telebot

# Create a non-root user
RUN addgroup --gid 1000 telebot
RUN adduser --ingroup telebot --shell /bin/sh telebot
USER telebot

# Command to run the application
CMD ["./telebot"]
