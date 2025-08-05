# Webex MCP Server Setup Guide

## 1. Build the Server

```bash
make build
```

## 2. Edit Claude Desktop Config

Find your config file:
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

Add this configuration:
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

## 3. Restart Claude Desktop

- Completely quit Claude (not just close window)
- Start Claude again

## 4. Use It!

In Claude, you can now say:
- "List my Webex rooms"
- "Send a message to the Engineering team"
- "Show recent messages in Project Alpha"
- "Create a new room called Planning"

## Important Notes

- The `command` path must be absolute (e.g., `/Users/john/webex-mcp-server-go/build/webex-mcp-server`)
- Get your Webex API token from [developer.webex.com](https://developer.webex.com)
- Make sure the binary is executable: `chmod +x /path/to/webex-mcp-server`

## Troubleshooting

If Claude doesn't see the Webex tools:
1. Check your JSON syntax is valid
2. Ensure the path to the binary is correct and absolute
3. Verify your Webex API token is valid
4. Check Claude Desktop developer tools (View > Toggle Developer Tools) for errors