# Builder stage
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# CGO_ENABLED=0 disables CGO for static linking
# GOOS=linux GOARCH=amd64 specifies the target OS and architecture
# -ldflags="-s -w" reduces the binary size by stripping debug info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o fiber-app .

# Final stage
FROM alpine:3.19

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/fiber-app /app/fiber-app

# Ensure log directory exists
RUN mkdir -p /app/logs && chown appuser:appuser /app/logs

USER appuser

EXPOSE 3000

ENTRYPOINT ["/app/fiber-app"]