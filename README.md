# Webex MCP Server (Go)

A Go implementation of the Webex Messaging MCP (Model Context Protocol) Server, providing AI assistants with full access to Cisco Webex messaging capabilities.

## Quick Setup for Claude Desktop

### 1. Build the Server
```bash
make build
```

### 2. Edit Claude Desktop Config
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

### 3. Restart Claude Desktop
- Completely quit Claude (not just close window)
- Start Claude again

### 4. Use It!
In Claude, you can now say:
- "List my Webex rooms"
- "Send a message to the Engineering team"
- "Show recent messages in Project Alpha"
- "Create a new room called Planning"

## Features

- Complete Webex API integration for messaging, rooms, teams, people, webhooks, and enterprise features
- MCP protocol support using the official Go SDK
- Support for both stdio and SSE (Server-Sent Events) transports
- Comprehensive tool registry with 53 Webex operations
- **Flexible HTTP client** with configurable backend:
  - Standard `net/http` for compatibility (default)
  - Optional `fasthttp` for 10x performance improvement
  - Automatic connection pooling and reuse
  - Configurable via `USE_FASTHTTP` environment variable

## Prerequisites

- Go 1.22 or higher
- Webex API token (obtain from [developer.webex.com](https://developer.webex.com))

## Installation

1. Clone the repository:
```bash
git clone https://github.com/raja-aiml/webex-mcp-server-go.git
cd webex-mcp-server-go
```

2. Install dependencies:
```bash
make deps
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env and add your Webex API token
```

## Configuration

Create a `.env` file with the following variables:

```env
WEBEX_PUBLIC_WORKSPACE_API_KEY=your_webex_api_token_here
WEBEX_API_BASE_URL=https://webexapis.com/v1
PORT=3001
USE_FASTHTTP=false  # Set to true for performance optimization
```

## Quick Start with Claude Desktop

**This MCP server is designed to work with Claude Desktop.** See the [Claude Desktop Quick Start Guide](docs/CLAUDE_DESKTOP_QUICK_START.md).

### 1. Build the server:
```bash
make build
```

### 2. Add to Claude Desktop config:
```json
{
  "mcpServers": {
    "webex": {
      "command": "/full/path/to/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "your-token"
      }
    }
  }
}
```

### 3. Restart Claude and use:
- "List my Webex rooms"
- "Send a message to [room name]"
- "Show recent messages"

## Other Usage Modes

### Run in stdio mode (for debugging):
```bash
./build/webex-mcp-server
```

### Run in HTTP mode (for health checks):
```bash
./build/webex-mcp-server -http :3001
curl http://localhost:3001/health
```

## Available Tools

The server provides the following tool categories:

### Messaging Tools
- `list_messages` - List messages in a room
- `create_a_message` - Send a message
- `delete_a_message` - Delete a message
- `edit_a_message` - Edit a message
- `get_message_details` - Get message details
- `list_direct_messages` - List direct messages

### Room Tools
- `list_rooms` - List rooms
- `create_a_room` - Create a room
- `get_room_details` - Get room details
- `update_a_room` - Update room settings
- `delete_a_room` - Delete a room
- `get_room_meeting_details` - Get room meeting info

### People Tools
- `list_people` - List people
- `create_a_person` - Create a person
- `get_person_details` - Get person details
- `update_a_person` - Update person info
- `delete_a_person` - Delete a person
- `get_my_own_details` - Get current user details

### Additional Tools
- Membership management
- Team and team membership operations
- Webhook management
- Room tabs
- Attachment actions
- Events
- ECM (Enterprise Content Management) operations

## Testing

**Important**: This MCP server is designed for use with Claude Desktop. The web-based MCP Inspector is NOT compatible.

### Quick Test
```bash
# Test stdio mode
./scripts/test_stdio.sh

# Test health endpoint (HTTP mode)
make run http  # Terminal 1
curl http://localhost:3001/health  # Terminal 2
```

### Claude Desktop Integration
Add to your Claude Desktop config to use all features:
- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`
- Linux: `~/.config/Claude/claude_desktop_config.json`

See [Testing Guide](docs/TESTING_CORRECTED.md) for details.

## Development

### Run tests:
```bash
make test
```

### Format code:
```bash
make fmt
```

### Clean build artifacts:
```bash
make clean
```

### Build for multiple platforms:
```bash
make build-all
```

## Architecture

The server follows SOLID principles with a clean, modular structure:

```
webex-mcp-server-go/
├── main.go                 - MCP server implementation
├── config/                 - Configuration management
│   └── config.go          - Environment and API configuration
├── webex/                  - Webex API client implementation
│   ├── client.go          - Unified HTTP client (supports net/http and fasthttp)
│   └── interface.go       - HTTPClient interface (DI/IoC)
└── tools/                  - MCP tool implementations (53 tools)
    ├── base.go            - Base tool functionality (DRY)
    ├── registry.go        - Tool registry and management
    ├── messages.go        - Message tools (6 tools)
    ├── people.go          - People tools (6 tools)
    ├── rooms.go           - Room tools (6 tools)
    ├── memberships.go     - Membership tools (5 tools)
    ├── teams.go           - Team tools (5 tools)
    ├── team_memberships.go - Team membership tools (5 tools)
    ├── webhooks.go        - Webhook tools (5 tools)
    ├── room_tabs.go       - Room tab tools (5 tools)
    ├── attachments.go     - Attachment action tools (2 tools)
    ├── events.go          - Event tools (2 tools)
    └── ecm.go             - ECM folder tools (5 tools)
```

### Design Principles Applied:
- **KISS**: Each file contains related tools only
- **YAGNI**: No over-engineering, just what's needed
- **SOLID**: Interface segregation, dependency inversion
- **DRY**: Shared base functionality, no code duplication



## License

MIT