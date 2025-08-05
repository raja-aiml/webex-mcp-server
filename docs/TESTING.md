# Testing the Webex MCP Server

This guide covers various methods to test the Webex MCP Server.

## Prerequisites

1. **Webex API Token**: Get your token from [developer.webex.com](https://developer.webex.com)
2. **Environment Setup**:
   ```bash
   cp .env.example .env
   # Edit .env and add your WEBEX_PUBLIC_WORKSPACE_API_KEY
   ```

## Testing Methods

### 1. Manual Testing with MCP Inspector

The MCP Inspector is a web-based tool for testing MCP servers.

```bash
# Start the server in HTTP mode
make run http

# Or directly:
./build/webex-mcp-server -http :3001

# The server will be available at http://localhost:3001
```

Then open the [MCP Inspector](https://modelcontextprotocol.io/inspector) and connect to `http://localhost:3001`.

### 2. Testing with Claude Desktop

Add the server to your Claude Desktop configuration:

**For macOS:**
Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "webex": {
      "command": "/path/to/webex-mcp-server-go/build/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

**For Windows:**
Edit `%APPDATA%\Claude\claude_desktop_config.json`

**For Linux:**
Edit `~/.config/Claude/claude_desktop_config.json`

Then restart Claude Desktop and you should see the Webex tools available.

### 3. Direct stdio Testing

Test the server directly using stdio mode:

```bash
# Build first
make build

# Run in stdio mode and send test commands
echo '{"jsonrpc": "2.0", "method": "tools/list", "id": 1}' | ./build/webex-mcp-server
```

### 4. HTTP API Testing

Test specific endpoints when running in HTTP mode:

```bash
# Start server
make run http

# Test health endpoint
curl http://localhost:3001/health

# Test with MCP protocol (using Server-Sent Events)
curl -X POST http://localhost:3001/sse \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "method": "tools/list", "id": 1}'
```

### 5. Automated Testing Script

Create a test script to verify basic functionality:

```bash
#!/bin/bash
# File: scripts/test_mcp_server.sh

# Start server in background
./build/webex-mcp-server -http :3001 &
SERVER_PID=$!

# Wait for server to start
sleep 2

# Test health endpoint
echo "Testing health endpoint..."
curl -f http://localhost:3001/health || exit 1

# Test tools list
echo -e "\n\nTesting tools list..."
curl -X POST http://localhost:3001/sse \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "method": "tools/list", "id": 1}'

# Clean up
kill $SERVER_PID
```

### 6. Testing Individual Tools

Test specific Webex tools through the MCP protocol:

```javascript
// Example: List rooms
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "list_rooms",
    "arguments": {
      "max": 10,
      "type": "group"
    }
  },
  "id": 2
}

// Example: Send a message
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "create_a_message",
    "arguments": {
      "roomId": "your-room-id",
      "text": "Hello from MCP!"
    }
  },
  "id": 3
}
```

### 7. Performance Testing

Use the included benchmark script:

```bash
# Run performance benchmarks
./scripts/benchmark.sh

# Or test with Apache Bench
ab -n 1000 -c 10 http://localhost:3001/health
```

### 8. Docker Testing

Test using Docker:

```bash
# Build and run with Docker
make docker build
make docker run

# Test with docker-compose
make docker run-dev
```

## Common Test Scenarios

### 1. List Available Tools
Verify all 53 tools are loaded:
```bash
echo '{"jsonrpc": "2.0", "method": "tools/list", "id": 1}' | ./build/webex-mcp-server | jq '.result | length'
# Should output: 53
```

### 2. Test Message Operations
```bash
# List messages in a room
# Create a test message
# Edit the message
# Delete the message
```

### 3. Test Room Management
```bash
# List rooms
# Create a room
# Update room details
# Add members
# Delete room
```

## Debugging

### Enable Debug Logging

Set environment variables for verbose logging:
```bash
export MCP_DEBUG=true
export LOG_LEVEL=debug
./build/webex-mcp-server
```

### Check Server Logs

When running in stdio mode, logs are written to stderr:
```bash
./build/webex-mcp-server 2>server.log
```

### Common Issues

1. **Authentication Errors**: Verify your API token is correct
2. **Connection Refused**: Check if the port is already in use
3. **Tool Not Found**: Ensure all tools are properly registered

## Integration Testing

For integration testing with actual Webex APIs:

1. Create a test room in Webex
2. Set up test data (users, messages, etc.)
3. Run through CRUD operations
4. Verify results in Webex client
5. Clean up test data

## Unit Testing

Run Go unit tests:
```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/tools/...
```

## Load Testing

For load testing the server:

```bash
# Using hey (HTTP load generator)
hey -n 10000 -c 100 http://localhost:3001/health

# Using vegeta
echo "GET http://localhost:3001/health" | vegeta attack -duration=30s -rate=100 | vegeta report
```

## Continuous Testing

Set up a watch script for development:
```bash
# Install entr or similar file watcher
ls **/*.go | entr -r make run http
```

This will restart the server whenever Go files change.