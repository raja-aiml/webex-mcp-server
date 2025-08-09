# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application for current platform
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/webex-mcp-server .

# Runtime stage - minimal image for MCP
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Create non-root user for security
RUN adduser -D -g '' mcpuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/build/webex-mcp-server /app/webex-mcp-server

# Make binary executable
RUN chmod +x /app/webex-mcp-server

# Create directory for .env file (optional - will be provided at runtime)
RUN mkdir -p /app/config

# Change ownership to non-root user
RUN chown -R mcpuser:mcpuser /app

# Switch to non-root user
USER mcpuser

# Set environment variable to look for .env in config directory (optional)
# The app will use environment variables if .env is not found
ENV DOTENV_PATH=/app/config/.env

# MCP servers typically run in stdio mode by default
ENTRYPOINT ["/app/webex-mcp-server"]

# Default to stdio mode for MCP
CMD []