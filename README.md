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
- **Optimized HTTP client**:
  - Standard `net/http` with connection pooling
  - Automatic connection reuse
  - Configurable timeouts and retry logic

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
```

## Quick Start with Claude Desktop

**This MCP server is designed to work with Claude Desktop.**

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
# Run all tests
make test

# Test specific components
make test-verbose   # Test with verbose output
make test-coverage  # Test with coverage report
```

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
│   ├── client.go          - HTTP client with connection pooling
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


# Docker Setup for Webex MCP Server

## Quick Start

```bash
# Test with MCP Inspector
source .env && npx @modelcontextprotocol/inspector go run .

# Start HTTP/SSE server
make run http

# Run in stdio mode
make run
```

### Manual Commands

```bash
# Build Docker image
docker build -t webex-mcp-server:latest .

# Run in stdio mode
docker run --rm -i --env-file .env webex-mcp-server:latest

# Run in HTTP mode
docker run --rm -p 8084:8084 --env-file .env webex-mcp-server:latest -http :8084

# Using docker-compose
docker-compose --profile http up --build
```

## Testing with MCP Inspector

### Option 1: Stdio Mode with Docker
```bash
npx @modelcontextprotocol/inspector docker run --rm -i --env-file .env webex-mcp-server:latest
```

### Option 2: HTTP/SSE Mode
1. Start the server:
```bash
docker-compose --profile http up -d
```

2. Open MCP Inspector and connect to:
```
http://localhost:8084/sse
```

## Docker Compose Profiles

- `stdio` - Run in stdio mode for MCP clients
- `http` - Run in HTTP/SSE mode on port 8084
- `dev` - Development mode with hot reload

## Troubleshooting

### .env File Not Found
The Docker container needs access to your `.env` file. The docker-compose.yml mounts it as a volume at `/app/config/.env`.

### Port Already in Use
Change the port in your `.env` file:
```
MCP_PORT=8085
```

### View Logs
```bash
docker-compose logs -f webex-mcp-http
```

### Stop Services
```bash
docker-compose --profile http down
```

## License

MIT