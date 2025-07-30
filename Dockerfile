# Multi-stage build for minimal production image
# Build stage
FROM golang:1.24.4-alpine AS builder

# Install build dependencies and security updates
RUN apk add --no-cache git ca-certificates tzdata && \
    apk upgrade --no-cache

# Create non-root user for build stage
RUN addgroup -g 1001 -S builder && \
    adduser -u 1001 -S builder -G builder

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY --chown=builder:builder go.mod go.sum ./

# Download and verify dependencies
RUN go mod download && \
    go mod verify

# Copy source code
COPY --chown=builder:builder *.go ./

# Switch to non-root user for building
USER builder

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags="-w -s -X main.version=1.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o cardgame-api .

# Runtime stage - using hardened Alpine
FROM alpine:3.19

# Install only essential runtime dependencies and apply security updates
RUN apk add --no-cache ca-certificates tzdata curl && \
    apk upgrade --no-cache && \
    rm -rf /var/cache/apk/*

# Create non-root user with specific UID/GID
RUN addgroup -g 65532 -S cardgame && \
    adduser -u 65532 -S cardgame -G cardgame -h /nonexistent -s /sbin/nologin

# Create app directory with proper permissions
RUN mkdir -p /app && \
    chown -R cardgame:cardgame /app

# Copy the binary with proper permissions
COPY --from=builder --chown=cardgame:cardgame /app/cardgame-api /usr/local/bin/cardgame-api

# Make binary executable
RUN chmod 755 /usr/local/bin/cardgame-api

# Copy static files and documentation
COPY --chown=cardgame:cardgame static/ /app/static/
COPY --chown=cardgame:cardgame api-docs.html /app/api-docs.html
COPY --chown=cardgame:cardgame openapi.yaml /app/openapi.yaml

# Set proper permissions on static files (read-only)
RUN find /app -type f -exec chmod 444 {} \; && \
    find /app -type d -exec chmod 555 {} \;

# Set up read-only root filesystem support
RUN mkdir -p /tmp && \
    chown cardgame:cardgame /tmp && \
    chmod 1777 /tmp

# Set working directory
WORKDIR /app

# Switch to non-root user
USER 65532:65532

# Set environment variables
ENV GIN_MODE=release
ENV LOG_LEVEL=INFO
ENV PORT=8080
ENV TRUSTED_PROXIES=""

# Expose port
EXPOSE 8080

# Add security labels
LABEL security.scan="true" \
      security.nonroot="true" \
      security.updates="auto"

# Health check using curl (more secure than wget)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/hello || exit 1

# Run the binary
ENTRYPOINT ["cardgame-api"]