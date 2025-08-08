# Instructions to Use Webex MCP Server

## Project Overview
This is a Go implementation of a Model Context Protocol (MCP) server that provides AI assistants with access to Cisco Webex messaging capabilities. The server integrates with Webex APIs and follows the MCP protocol using the official Go SDK.

## Configuration

**Important: Always load environment variables from `.env` file before running the server:**
```bash
source .env
```

Required environment variables in `.env` file:
```env
# Main API key for Webex workspace operations
WEBEX_PUBLIC_WORKSPACE_API_KEY=<your_workspace_api_key>
# Webex API base URL (default: https://webexapis.com/v1)
WEBEX_API_BASE_URL=https://webexapis.com/v1
# Server port for HTTP mode (default: 3001)
PORT=3001
# Optional: Bot access token for bot-specific operations
WEBEX_BOT_ACCESS_TOKEN=<your_bot_token>
```

Note: Never commit API keys to version control. Keep `.env` in `.gitignore`.

## Running Modes

1. **VSCode Extension Mode** (Default):
   ```bash
   ./build/webex-mcp-server
   ```
   Configure in VSCode settings.json:
   ```json
   {
     "mcp.webex.server": {
       "path": "/path/to/webex-mcp-server",
       "env": {
         "WEBEX_PUBLIC_WORKSPACE_API_KEY": "${env:WEBEX_PUBLIC_WORKSPACE_API_KEY}"
       }
     }
   }
   ```

2. **HTTP Server Mode**:
   ```bash
   ./build/webex-mcp-server -http :3001
   ```

## Available Tools

### Rooms
- `list_rooms` - List Webex rooms
- `create_room` - Create a new room
- `get_room` - Get room details

### Messages
- `create_message` - Send a message
- `list_messages` - List messages in a room

### People
- `list_people` - List people
- `get_me` - Get own details

### Memberships
- `create_membership` - Add person to room
- `list_memberships` - List room memberships

### Teams
- `create_team` - Create a new team
- `list_teams` - List teams

### Webhooks
- `create_webhook` - Create a webhook
- `list_webhooks` - List webhooks

## Example Workflow: Create Room and Send Message

1. Initialize session
2. Create a new room
3. Send a message to the room
4. Add team members (optional)

Each operation requires proper JSON-RPC 2.0 format with `jsonrpc: "2.0"` and `id` fields.

## Important Notes
1. Always initialize a session before making tool calls
2. HTTP mode requires `Accept: application/json, text/event-stream` header
3. Replace `ROOM_ID` placeholders with actual room IDs
4. Some operations may require additional permissions

## Key Architecture Components

### Core Components
- `main.go` - Entry point, handles CLI flags and app initialization
- `internal/app/` - Application lifecycle management
- `internal/server/` - MCP server implementation with stdio/HTTP transport
- `internal/webex/` - Webex API client with flexible HTTP backend
- `internal/tools/` - Tool implementations (53 Webex operations)

### Design Patterns
1. **Modular Architecture**
   - Each tool category (messages, rooms, people, etc.) has its own file
   - Base functionality shared through `tools/base.go`
   - Tool registry pattern in `tools/registry.go`

2. **Flexible HTTP Client**
   - Supports both `net/http` and `fasthttp` backends
   - Configured via `USE_FASTHTTP` environment variable
   - Interface defined in `webex/interface.go`

## Development Workflows

### Build and Test
```bash
make build           # Build the server
make deps            # Install dependencies
make test           # Run all tests
make fmt            # Format code
./scripts/test.sh   # Run specific test suites
```

### Configuration

**Important: Always load environment variables from `.env` file before running the server:**
```bash
source .env
```

Required environment variables should be placed in `.env` file:
```env
# Main API key for Webex workspace operations
WEBEX_PUBLIC_WORKSPACE_API_KEY=<your_workspace_api_key>

# Webex API base URL (default: https://webexapis.com/v1)
WEBEX_API_BASE_URL=https://webexapis.com/v1

# Server port for HTTP mode (default: 3001)
PORT=3001

# Toggle FastHTTP client (default: false)
USE_FASTHTTP=false

# Optional: Bot access token for bot-specific operations
WEBEX_BOT_ACCESS_TOKEN=<your_bot_token>
```

Note: 
- Never commit actual API keys to version control
- Keep your `.env` file in `.gitignore`
- You can copy `.env.example` to `.env` and fill in your keys

### Running Modes

The server can be used in three different modes:

1. **VSCode Extension Mode**:
   ```bash
   # Start in stdio mode (default)
   ./build/webex-mcp-server
   ```
   Configure in VSCode settings.json:
   ```json
   {
     "mcp.webex.server": {
       "path": "/path/to/webex-mcp-server",
       "env": {
         "WEBEX_PUBLIC_WORKSPACE_API_KEY": "${env:WEBEX_PUBLIC_WORKSPACE_API_KEY}"
       }
     }
   }
   ```

2. **HTTP Server Mode**:
   ```bash
   # Start HTTP server on port 3001
   ./build/webex-mcp-server -http :3001
   ```
   Example API calls:
   ```bash
   # Initialize session
   curl -X POST http://localhost:3001/mcp/initialize \
     -H "Content-Type: application/json" \
     -H "Accept: application/json, text/event-stream" \
     -d '{
       "jsonrpc": "2.0",
       "id": 1,
       "method": "initialize",
       "params": {
         "name": "test-client",
         "version": "1.0.0",
         "capabilities": {}
       }
     }'

   # Call list_rooms tool (after initialization)
   curl -X POST http://localhost:3001/mcp/tool/call \
     -H "Content-Type: application/json" \
     -H "Accept: application/json, text/event-stream" \
     -d '{
       "jsonrpc": "2.0",
       "id": 2,
       "method": "tool/call",
       "params": {
         "name": "list_rooms",
         "arguments": {
           "max": 10,
           "type": "group"
         }
       }
     }'
   ```

3. **Claude Desktop Mode**:
   ```bash
   # Start in stdio mode (default)
   ./build/webex-mcp-server
   ```
   MCP Protocol Example:
   ```json
   // Initialize session
   {
     "jsonrpc": "2.0",
     "id": 1,
     "method": "initialize",
     "params": {
       "name": "claude-desktop",
       "version": "1.0.0",
       "capabilities": {}
     }
   }

   // Call tool
   {
     "jsonrpc": "2.0",
     "id": 2,
     "method": "tool/call",
     "params": {
       "name": "list_rooms",
       "arguments": {
         "max": 10,
         "type": "group"
       }
     }
   }
   ```

Important Notes:
1. For all modes, ensure environment variables are properly set in `.env` file or system environment.
2. The server implements JSON-RPC 2.0 protocol - all requests must include `jsonrpc: "2.0"` and `id` fields.
3. HTTP mode requires `Accept: application/json, text/event-stream` header for proper response handling.
4. Tool calls must be made after a successful session initialization.
5. Responses are returned as Server-Sent Events (SSE) in HTTP mode.

### Example Workflow: Create Room and Send Message
Here's a complete workflow showing how to create a room and send a message:

1. Initialize session:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "name": "test-client",
    "version": "1.0.0",
    "capabilities": {}
  }
}
```

2. Create a new room:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tool/call",
  "params": {
    "name": "create_room",
    "arguments": {
      "title": "Project Discussion",
      "type": "group"
    }
  }
}
```

3. Send a message to the new room (using roomId from previous response):
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tool/call",
  "params": {
    "name": "create_message",
    "arguments": {
      "roomId": "ROOM_ID_FROM_PREVIOUS_RESPONSE",
      "text": "Welcome to the Project Discussion room!"
    }
  }
}
```

4. Add a team member (optional):
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tool/call",
  "params": {
    "name": "create_membership",
    "arguments": {
      "roomId": "ROOM_ID_FROM_PREVIOUS_RESPONSE",
      "personEmail": "colleague@example.com"
    }
  }
}
```

## Available Tools and Examples

### Rooms
- `list_rooms` - List Webex rooms
  ```json
  {
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tool/call",
    "params": {
      "name": "list_rooms",
      "arguments": {
        "max": 10,
        "type": "group"
      }
    }
  }
  ```

- `create_room` - Create a new room
  ```json
  {
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tool/call",
    "params": {
      "name": "create_room",
      "arguments": {
        "title": "My New Room",
        "type": "group"
      }
    }
  }
  ```

- `get_room` - Get room details
  ```json
  {
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tool/call",
    "params": {
      "name": "get_room",
      "arguments": {
        "roomId": "ROOM_ID"
      }
    }
  }
  ```

### Messages
- `create_message` - Send a message
  ```json
  {
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tool/call",
    "params": {
      "name": "create_message",
      "arguments": {
        "roomId": "ROOM_ID",
        "text": "Hello, World!"
      }
    }
  }
  ```

- `list_messages` - List messages in a room
  ```json
  {
    "jsonrpc": "2.0",
    "id": 5,
    "method": "tool/call",
    "params": {
      "name": "list_messages",
      "arguments": {
        "roomId": "ROOM_ID",
        "max": 50
      }
    }
  }
  ```

### People
- `list_people` - List people
  ```json
  {
    "jsonrpc": "2.0",
    "id": 6,
    "method": "tool/call",
    "params": {
      "name": "list_people",
      "arguments": {
        "email": "example@webex.com",
        "max": 10
      }
    }
  }
  ```

- `get_me` - Get own details
  ```json
  {
    "jsonrpc": "2.0",
    "id": 7,
    "method": "tool/call",
    "params": {
      "name": "get_me",
      "arguments": {}
    }
  }
  ```

### Memberships
- `create_membership` - Add person to room
  ```json
  {
    "jsonrpc": "2.0",
    "id": 8,
    "method": "tool/call",
    "params": {
      "name": "create_membership",
      "arguments": {
        "roomId": "ROOM_ID",
        "personEmail": "person@example.com"
      }
    }
  }
  ```

- `list_memberships` - List room memberships
  ```json
  {
    "jsonrpc": "2.0",
    "id": 9,
    "method": "tool/call",
    "params": {
      "name": "list_memberships",
      "arguments": {
        "roomId": "ROOM_ID",
        "max": 100
      }
    }
  }
  ```

### Teams
- `create_team` - Create a new team
  ```json
  {
    "jsonrpc": "2.0",
    "id": 10,
    "method": "tool/call",
    "params": {
      "name": "create_team",
      "arguments": {
        "name": "My New Team"
      }
    }
  }
  ```

- `list_teams` - List teams
  ```json
  {
    "jsonrpc": "2.0",
    "id": 11,
    "method": "tool/call",
    "params": {
      "name": "list_teams",
      "arguments": {
        "max": 100
      }
    }
  }
  ```

### Webhooks
- `create_webhook` - Create a webhook
  ```json
  {
    "jsonrpc": "2.0",
    "id": 12,
    "method": "tool/call",
    "params": {
      "name": "create_webhook",
      "arguments": {
        "name": "My Webhook",
        "targetUrl": "https://example.com/webhook",
        "resource": "messages",
        "event": "created"
      }
    }
  }
  ```

- `list_webhooks` - List webhooks
  ```json
  {
    "jsonrpc": "2.0",
    "id": 13,
    "method": "tool/call",
    "params": {
      "name": "list_webhooks",
      "arguments": {}
    }
  }
  ```

Important Notes:
1. Replace `ROOM_ID` with actual room IDs from your Webex workspace
2. All examples assume you've already initialized the session
3. HTTP mode requires proper headers as shown in the Running Modes section
4. Some operations may require additional permissions in your Webex token

## Key Integration Points

### 1. MCP Protocol Integration
- Tool registry in `tools/registry.go` defines available operations
- Each tool implements the MCP tool interface
- Example: `tools/messages.go` for message operations

### 2. Webex API Integration
- HTTP client in `webex/client.go`
- API endpoints defined per tool category
- Authentication via `WEBEX_PUBLIC_WORKSPACE_API_KEY`

## Common Development Tasks

### Adding New Tools
1. Create tool implementation in appropriate file under `tools/`
2. Register tool in `tools/registry.go`
3. Add tests in corresponding `_test.go` file

### Error Handling Pattern
- Use descriptive error types
- Wrap errors with context
- Return errors to be handled by caller

### Testing Strategy
- Unit tests per component
- Integration tests for API interactions
- Benchmark tests for performance validation

## Reference Examples

### Tool Implementation
```go
// From tools/messages.go
type ListMessagesRequest struct {
    RoomID string `json:"roomId"`
    // ... other fields
}
```

### HTTP Client Usage
```go
// From webex/client.go
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
    // ... implementation
}
```

## Common Pitfalls
1. **Configuration Management**
   - Always initialize config using `config.Load()` before using tools
   - Ensure all required environment variables (especially `WEBEX_PUBLIC_WORKSPACE_API_KEY`) are set
   - Use `config.MustLoad()` only in initialization code where errors are fatal

2. **Transport Mode Handling**
   - Support both stdio (default) and HTTP transport modes
   - Handle proper initialization and shutdown for each mode
   - Use correct port configuration for HTTP mode

3. **Resource Management**
   - Properly clean up resources in shutdown
   - Handle context cancellation in long-running operations
   - Use sync.Once for singleton initialization

4. **Input Validation**
   - Validate all tool inputs before making API calls
   - Convert and validate schema formats correctly
   - Handle both direct jsonschema.Schema and legacy schemas

5. **Error Handling**
   - Return tool errors as results per MCP spec, not as protocol errors
   - Use proper error wrapping with context
   - Handle both synchronous and asynchronous errors appropriately

6. **Security Considerations**
   - Clean and validate API tokens before use
   - Set appropriate HTTP headers for API requests
   - Handle authentication failures gracefully
