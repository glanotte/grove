# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates (needed for go modules)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grove .

# Final stage
FROM alpine:latest

# Install git, docker, and other necessary tools
RUN apk add --no-cache git docker docker-compose ca-certificates

# Create non-root user
RUN addgroup -g 1000 grove && \
    adduser -u 1000 -G grove -s /bin/sh -D grove

# Set working directory
WORKDIR /home/grove

# Copy binary from builder stage
COPY --from=builder /app/grove /usr/local/bin/grove

# Copy templates and examples
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/examples ./examples
COPY --from=builder /app/scripts ./scripts

# Make binary executable
RUN chmod +x /usr/local/bin/grove

# Switch to non-root user
USER grove

# Set up environment
ENV PATH="/usr/local/bin:${PATH}"

# Default command
CMD ["grove", "--help"]