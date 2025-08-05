# Webex MCP Server Setup

## Quick Setup

### 1. Build
```bash
make build
```

### 2. Configure Claude Desktop

Find your config file:
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

Add:
```json
{
  "mcpServers": {
    "webex": {
      "command": "/absolute/path/to/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "your-webex-token"
      }
    }
  }
}
```

### 3. Restart Claude Desktop
Quit completely and restart.

### 4. Use
- "List my Webex rooms"
- "Send a message to [room name]"
- "Show recent messages"

## Troubleshooting

**Not working?**
1. Path must be absolute
2. Check JSON syntax: `python3 -m json.tool ~/Library/Application\ Support/Claude/claude_desktop_config.json`
3. Make executable: `chmod +x /path/to/webex-mcp-server`
4. Test: `./build/webex-mcp-server` (should show stdio mode)
5. Check Claude Developer Tools (View > Toggle Developer Tools)

**Get Webex Token**
1. Visit [developer.webex.com](https://developer.webex.com)
2. Sign in → My Webex Apps → Personal Access Token

## Testing

```bash
# Test basic functionality
./scripts/test.sh

# Test specific component
./scripts/test.sh health  # Health endpoint
./scripts/test.sh stdio   # Stdio mode
./scripts/test.sh bench   # Performance