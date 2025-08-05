# Setting Up Webex MCP Server with Claude Desktop

This guide walks you through setting up the Webex MCP Server to work with Claude Desktop.

## Prerequisites

1. **Claude Desktop** - Download from [claude.ai](https://claude.ai/download)
2. **Webex API Token** - Get from [developer.webex.com](https://developer.webex.com)
3. **Built MCP Server** - Follow the build instructions in the main README

## Step 1: Build the Server

```bash
# Clone the repository
git clone https://github.com/raja-aiml/webex-mcp-server-go.git
cd webex-mcp-server-go

# Build the server
make build

# Note the full path to the binary
pwd
# Example output: /Users/yourname/webex-mcp-server-go
# Your binary will be at: /Users/yourname/webex-mcp-server-go/build/webex-mcp-server
```

## Step 2: Locate Claude Desktop Configuration

Find your Claude Desktop configuration file:

### macOS
```bash
# Open the config directory
open ~/Library/Application\ Support/Claude/

# The config file is:
~/Library/Application Support/Claude/claude_desktop_config.json
```

### Windows
```cmd
# Open File Explorer to:
%APPDATA%\Claude\

# The config file is:
%APPDATA%\Claude\claude_desktop_config.json
```

### Linux
```bash
# The config file is:
~/.config/Claude/claude_desktop_config.json
```

## Step 3: Edit Configuration

1. Open the `claude_desktop_config.json` file in a text editor
2. Add the Webex MCP server configuration:

```json
{
  "mcpServers": {
    "webex": {
      "command": "/Users/yourname/webex-mcp-server-go/build/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "your-actual-webex-token-here"
      }
    }
  }
}
```

**Important**: 
- Replace `/Users/yourname/webex-mcp-server-go/build/webex-mcp-server` with YOUR actual full path
- Replace `your-actual-webex-token-here` with your Webex API token
- Use forward slashes (/) even on Windows

### Example Complete Configuration

If you have other MCP servers, your config might look like:

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-filesystem", "/Users/yourname/Documents"]
    },
    "webex": {
      "command": "/Users/yourname/webex-mcp-server-go/build/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "MDY3ZmVjNTgtNGI3Yi00NDQ3LWEwNDMtZTY5ZDA2MzViOTg3YmE5NjM0MTUtYzE2_PF84_1eb65fdf-9643-417f-9974-ad72cae0e10f"
      }
    }
  }
}
```

## Step 4: Restart Claude Desktop

1. Completely quit Claude Desktop (not just close the window)
   - macOS: Cmd+Q or Claude > Quit Claude
   - Windows: Right-click system tray icon > Exit
   - Linux: Close and ensure process is terminated

2. Start Claude Desktop again

## Step 5: Verify Connection

1. Open Claude Desktop
2. Look for the MCP server indicator (usually shows connected servers)
3. Start a new conversation
4. Type: "What Webex tools do you have available?"

Claude should respond with a list of 53 Webex tools.

## Step 6: Using Webex with Claude

Here are example prompts you can use:

### Basic Commands
```
"List my Webex rooms"
"Show my recent Webex messages"
"What's my Webex user information?"
"List people in my Webex organization"
```

### Messaging
```
"Send a Webex message saying 'Hello team!' to the room called 'Project Updates'"
"Create a new Webex room called 'Planning Meeting'"
"Add john@example.com to the 'Planning Meeting' room"
```

### Advanced Usage
```
"Show me the last 10 messages in my 'General' Webex room"
"List all Webex rooms I'm a member of and show their last activity"
"Find all Webex users whose name contains 'Smith'"
```

## Troubleshooting

### Claude doesn't see Webex tools

1. **Check the config file syntax** - Make sure it's valid JSON
   ```bash
   # Validate JSON (macOS/Linux)
   python3 -m json.tool ~/Library/Application\ Support/Claude/claude_desktop_config.json
   ```

2. **Verify the binary path** - Make sure the path is absolute and correct
   ```bash
   # Test if the binary exists and is executable
   ls -la /path/to/your/webex-mcp-server
   ```

3. **Check the server starts correctly**
   ```bash
   # Test standalone
   /path/to/your/webex-mcp-server
   # Should see: Starting webex-mcp-server v0.1.0 in stdio mode
   # Press Ctrl+C to exit
   ```

4. **Verify API token** - Make sure your Webex token is valid
   ```bash
   # Test API token
   curl -H "Authorization: Bearer YOUR_TOKEN" https://webexapis.com/v1/people/me
   ```

### Common Issues

1. **"command not found"** - Use absolute path, not relative
2. **"permission denied"** - Make sure binary is executable: `chmod +x /path/to/webex-mcp-server`
3. **No Webex data** - Check your API token is correct and has proper scopes
4. **Server crashes** - Check Claude Desktop logs:
   - macOS: `~/Library/Logs/Claude/`
   - Windows: `%LOCALAPPDATA%\Claude\logs\`

### Viewing Logs

To see what's happening:

1. **Claude Desktop Developer Tools**
   - In Claude Desktop: View > Toggle Developer Tools
   - Check Console for errors

2. **MCP Server Logs**
   - Servers log to stderr, which Claude captures
   - Check Claude's developer console for server output

## Security Notes

1. **API Token Security**
   - Your Webex API token is stored in plain text in the config
   - Ensure your config file has appropriate permissions
   - Consider using a token with limited scopes

2. **File Permissions**
   ```bash
   # Restrict config file access (macOS/Linux)
   chmod 600 ~/Library/Application\ Support/Claude/claude_desktop_config.json
   ```

## Next Steps

Now that Webex is connected to Claude, you can:

1. Ask Claude to manage your Webex communications
2. Create automations and workflows
3. Analyze message patterns
4. Generate reports from Webex data

Remember: Claude can see and interact with your Webex data based on the permissions of your API token.