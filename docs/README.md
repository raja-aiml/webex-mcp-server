# Webex MCP Server Documentation

## Setup Guides

- **[Quick Setup](SETUP.md)** - Simple 4-step setup guide
- **[Claude Desktop Quick Start](CLAUDE_DESKTOP_QUICK_START.md)** - Visual quick start guide with examples
- **[Claude Desktop Detailed Setup](CLAUDE_DESKTOP_SETUP.md)** - Comprehensive setup with troubleshooting

## Technical Documentation

- **[Project Structure](PROJECT_STRUCTURE.md)** - Codebase organization and architecture
- **[Testing Guide](TESTING_CORRECTED.md)** - How to test the MCP server

## Quick Links

### Setup in 30 Seconds
1. `make build`
2. Add to Claude config: `/path/to/webex-mcp-server`
3. Restart Claude
4. Say "List my Webex rooms"

### Config File Locations
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

### Get Help
- Check the [Troubleshooting section](CLAUDE_DESKTOP_SETUP.md#troubleshooting)
- View Claude Developer Tools: View > Toggle Developer Tools
- Test server: `./build/webex-mcp-server` (should show stdio mode message)