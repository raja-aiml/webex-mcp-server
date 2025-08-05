# Testing the Webex MCP Server - Corrected Guide

## Important Note About MCP Inspector

The web-based MCP Inspector at https://modelcontextprotocol.io/inspector is **NOT compatible** with this Go-based MCP server. The Go SDK uses StreamableHTTP transport which is designed for:
- Claude Desktop integration
- Programmatic clients
- Direct stdio communication

## Correct Testing Methods

### 1. Testing with Claude Desktop (Recommended)

This is the primary way to use MCP servers. Add to your Claude Desktop configuration:

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows:** `%APPDATA%\Claude\claude_desktop_config.json`
**Linux:** `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "webex": {
      "command": "/full/path/to/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

### 2. Direct stdio Testing

The server is designed to communicate via stdio:

```bash
# Test tool listing
echo '{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}, "id": 1}' | ./build/webex-mcp-server

# The server expects proper MCP protocol initialization
```

### 3. Using the Local MCP Inspector (if needed)

Install and run the MCP Inspector locally with Node.js:

```bash
# Install the inspector
npm install -g @modelcontextprotocol/inspector

# Run with stdio server
npx @modelcontextprotocol/inspector ./build/webex-mcp-server
```

### 4. HTTP Mode Testing

When running in HTTP mode (`-http` flag), the server provides:
- `/health` - Health check endpoint
- `/` - StreamableHTTP handler (not REST API)

```bash
# Start server
./build/webex-mcp-server -http :3001

# Test health
curl http://localhost:3001/health
```

### 5. Testing Tools Programmatically

The best way to test individual tools is through Claude Desktop or by implementing a proper MCP client that supports the StreamableHTTP transport.

## Why the Web Inspector Doesn't Work

The web-based MCP Inspector expects servers to implement specific transport protocols that differ from the Go SDK's implementation. The Go SDK uses:
- StreamableHTTP for HTTP mode
- stdio for command mode

These are designed for production use with Claude Desktop and other MCP clients, not for web-based debugging tools.

## Recommended Testing Approach

1. **Use Claude Desktop** - This is how end users will interact with your server
2. **Write unit tests** - Test your tool implementations directly
3. **Use logging** - Run with stderr output to see what's happening:
   ```bash
   ./build/webex-mcp-server 2>server.log
   ```

## Example: Testing with a Simple Client

Here's a minimal Python client that can communicate with the stdio server:

```python
import subprocess
import json

# Start the server
proc = subprocess.Popen(
    ['./build/webex-mcp-server'],
    stdin=subprocess.PIPE,
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE,
    text=True
)

# Send initialize request
request = {
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {
        "protocolVersion": "2024-11-05",
        "capabilities": {}
    },
    "id": 1
}

proc.stdin.write(json.dumps(request) + '\n')
proc.stdin.flush()

# Read response
response = proc.stdout.readline()
print("Response:", response)
```

## Summary

- ✅ Use Claude Desktop for real-world testing
- ✅ Use stdio mode for debugging
- ✅ Use health endpoint for monitoring
- ❌ Don't use web-based MCP Inspector
- ❌ Don't expect REST API endpoints