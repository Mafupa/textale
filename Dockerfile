# Build stage
FROM golang:1.21-bookworm AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y gcc sqlite3 libsqlite3-dev && rm -rf /var/lib/apt/lists/*

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=1 go build -o textale .

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /app/textale .

# Create directories
RUN mkdir -p .ssh data

# Expose SSH port
EXPOSE 2222

# Run the application
CMD ["./textale"]
