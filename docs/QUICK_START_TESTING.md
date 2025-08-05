# Quick Start: Testing the Webex MCP Server

## 1. Basic Setup

```bash
# Clone and build
git clone https://github.com/raja-aiml/webex-mcp-server-go.git
cd webex-mcp-server-go
make build

# Set up environment
cp .env.example .env
# Edit .env and add your WEBEX_PUBLIC_WORKSPACE_API_KEY
```

## 2. Quick Test - Automated

```bash
# Run automated tests
./scripts/test_mcp_server.sh
```

This will:
- Start the server
- Test health endpoint
- List all available tools
- Test error handling
- Show results

## 3. Interactive Testing

```bash
# Terminal 1: Start the server
make run http

# Terminal 2: Run interactive client
./scripts/mcp_client.sh
```

The interactive client lets you:
- List and explore all 53 tools
- Test Webex operations
- Send custom MCP requests

## 4. Test with Claude Desktop

1. Add to Claude Desktop config:
```json
{
  "mcpServers": {
    "webex": {
      "command": "/full/path/to/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "your-token-here"
      }
    }
  }
}
```

2. Restart Claude Desktop

3. In Claude, you can now say:
   - "List my Webex rooms"
   - "Send a message to [room name]"
   - "Show recent messages"

## 5. Manual stdio Test

```bash
# Quick tool list test
echo '{"jsonrpc": "2.0", "method": "tools/list", "id": 1}' | ./build/webex-mcp-server | jq '.result[].name'
```

## 6. HTTP API Test

```bash
# Health check
curl http://localhost:3001/health | jq .

# List tools via HTTP
curl -X POST http://localhost:3001 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "method": "tools/list", "id": 1}' | jq '.result | length'
```

## Common Test Commands

### List Your Rooms
```bash
curl -X POST http://localhost:3001 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "list_rooms",
      "arguments": {"max": 5}
    },
    "id": 1
  }' | jq .
```

### Get Your User Info
```bash
curl -X POST http://localhost:3001 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "get_my_own_details",
      "arguments": {}
    },
    "id": 1
  }' | jq .
```

## Troubleshooting

1. **Server won't start**: Check if port 3001 is already in use
2. **Authentication errors**: Verify your API token in .env
3. **No tools found**: Ensure the build completed successfully
4. **Connection refused**: Make sure server is running with `-http` flag

## Performance Test

```bash
# Quick performance check (requires Apache Bench)
ab -n 100 -c 10 http://localhost:3001/health
```

## Docker Test

```bash
# Test with Docker
make docker build
make docker run
```

Then test at http://localhost:8084/health