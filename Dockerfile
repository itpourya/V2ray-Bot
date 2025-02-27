# Stage 1: Build the Go application
FROM golang:1.20 AS builder

WORKDIR /app

# Copy go.mod and go.sum files for dependency resolution
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application as a static binary
RUN CGO_ENABLED=0 go build -o redzone ./cmd/main.go

# Stage 2: Create a minimal image for the application
FROM alpine:latest

# Install timezone data
RUN apk add --no-cache tzdata

# Set the timezone environment variable
ENV TZ=Asia/Tehran

# Create a non-root user and group
RUN addgroup --gid 1000 redzone \
    && adduser --disabled-password --ingroup redzone --shell /bin/sh redzone

# Set the working directory to a non-root directory
WORKDIR /home/redzone

# Copy the built application from the builder stage
COPY --from=builder /app/redzone .

# Ensure the binary is executable
RUN chmod +x ./redzone

# Switch to the non-root user
USER redzone

ENV TOKEN="7121293474:AAHPfmIZbIQ0ZsI7sgA0wMnxnZY33ymAy2s"

# Command to run the application
CMD ["./redzone"]
