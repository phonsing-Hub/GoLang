# Dockerfile
# Multi-stage build for smaller production image

#############
# Build Stage
#############
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser for security
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0: ปิด CGO เพื่อให้ binary สามารถรันใน scratch image ได้
# GOOS=linux: target OS
# -a: force rebuilding
# -installsuffix cgo: add suffix to package installation directory
# -ldflags: linker flags
#   -s: omit symbol table and debug info
#   -w: omit DWARF symbol table
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o app ./main.go

####################
# Production Stage
####################
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /build/app /app

# Copy static files if any (optional)
# COPY --from=builder /build/static /static
# COPY --from=builder /build/templates /templates

# Use non-root user
USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app", "healthcheck"] || exit 1

# Run the application
ENTRYPOINT ["/app"]