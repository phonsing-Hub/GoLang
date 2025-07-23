# --- Builder stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Use cache effectively
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build main app and migrate together (ใช้ cache เดียว)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o fiber-app .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o migrate ./scripts/migration.go

# --- Runner stage ---
FROM alpine:3.19

RUN adduser -D -g '' appuser
WORKDIR /app

COPY --from=builder /app/fiber-app .
COPY --from=builder /app/migrate .

RUN mkdir -p /app/logs && chown appuser:appuser /app/logs

USER appuser
EXPOSE 3000
ENTRYPOINT ["/app/fiber-app"]