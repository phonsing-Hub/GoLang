#############
# Build Stage
#############
FROM golang:1.24-alpine AS builder 

# install tools, workdir, copy go.mod, build binary
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./main.go

####################
# Production Stage
####################
FROM alpine:3.19

RUN adduser -D -g '' appuser && \
    apk add --no-cache tzdata

RUN mkdir -p /app/logs && chown appuser:appuser /app/logs

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app /app/app

USER appuser

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/app", "healthcheck"] || exit 1

ENTRYPOINT ["/app/app"]