# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o surfmate.io ./main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies for browser automation
RUN apk add --no-cache \
    chromium \
    chromium-chromedriver \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

# Copy binary from builder
COPY --from=builder /app/surfmate.io /usr/local/bin/surfmate.io

# Set environment variables for headless Chrome in container
ENV BROWSER_PATH=/usr/bin/chromium-browser
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROMIUM_FLAGS="--no-sandbox --disable-dev-shm-usage"

# Expose HTTP port (for HTTP mode)
EXPOSE 8080

# Run the application
ENTRYPOINT ["surfmate.io"]
