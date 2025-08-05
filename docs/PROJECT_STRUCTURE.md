# Project Structure

This project follows Go best practices with a clean, modular architecture:

```
webex-mcp-server-go/
├── cmd/
│   └── webex-mcp-server/
│       └── main.go                 # Application entry point
├── internal/                       # Private packages (not importable by other projects)
│   ├── app/
│   │   └── app.go                 # Application lifecycle management
│   ├── config/
│   │   ├── config.go              # Configuration management
│   │   └── provider.go            # Config provider interface
│   ├── handlers/
│   │   └── handlers.go            # HTTP handlers and health checks
│   ├── server/
│   │   ├── server.go              # MCP server creation and tool registration
│   │   ├── config.go              # Server configuration initialization
│   │   └── transport.go           # HTTP and stdio transport implementations
│   ├── tools/                     # MCP tools implementation
│   │   ├── registry.go            # Tool registry
│   │   ├── base.go                # Base tool functionality
│   │   ├── tool_factory.go        # Tool factory utilities
│   │   ├── generic_tool.go        # Generic tool implementation
│   │   ├── messages*.go           # Messaging tools
│   │   ├── rooms.go               # Room management tools
│   │   ├── people.go              # People management tools
│   │   ├── teams*.go              # Team management tools
│   │   ├── memberships.go         # Membership tools
│   │   ├── webhooks.go            # Webhook tools
│   │   ├── attachments.go         # Attachment tools
│   │   ├── events.go              # Event tools
│   │   ├── ecm.go                 # ECM tools
│   │   └── room_tabs.go           # Room tab tools
│   └── webex/                     # Webex API client
│       ├── client.go              # HTTP client implementation
│       └── interface.go           # Client interface definition
├── deployment/                     # Deployment configurations
│   ├── docker/
│   │   ├── Dockerfile             # Multi-stage Docker build
│   │   └── .dockerignore          # Docker ignore patterns
│   └── local/
│       └── docker-compose.yaml    # Local development orchestration
├── scripts/                        # Build and utility scripts
│   ├── release.sh                 # Release automation
│   ├── benchmark.sh               # Performance benchmarking
│   ├── test_server.sh             # Test server setup
│   └── update_imports.sh          # Import path migration
├── build/                          # Build artifacts (gitignored)
├── .env.example                    # Environment configuration template
├── .gitignore                      # Git ignore patterns
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
├── Makefile                        # Build and development tasks
└── README.md                       # Project documentation
```

## Package Descriptions

### cmd/webex-mcp-server
The main application entry point. Contains minimal code - just flag parsing and application initialization.

### internal/app
Application lifecycle management including:
- Configuration initialization
- Server creation and startup
- Signal handling for graceful shutdown
- Context management

### internal/config
Configuration management with:
- Environment variable loading
- Configuration validation
- Provider interface for dependency injection

### internal/handlers
HTTP request handlers:
- Health check endpoint
- MCP protocol handlers
- HTTP multiplexer setup

### internal/server
MCP server implementation:
- Server creation and configuration
- Tool registration
- Transport layer (HTTP/SSE and stdio)

### internal/tools
Comprehensive MCP tool implementations:
- 53 Webex API operations
- Generic tool framework
- Tool registry and factory patterns
- Base functionality for DRY principles

### internal/webex
Webex API client:
- Configurable HTTP backend (net/http or fasthttp)
- Automatic connection pooling
- Error handling
- Interface for dependency injection

## Design Principles

1. **Standard Go Project Layout**: Follows community conventions with `cmd/` and `internal/` directories
2. **Dependency Injection**: Interfaces and providers for testability
3. **SOLID Principles**: Single responsibility, open/closed, interface segregation
4. **DRY (Don't Repeat Yourself)**: Shared base functionality, generic implementations
5. **Clean Architecture**: Clear separation of concerns between packages
6. **12-Factor App**: Configuration through environment, stateless processes

## Build and Run

```bash
# Build the application
make build

# Run in stdio mode
make run

# Run in HTTP/SSE mode
make run http

# Run with Docker
make docker build
make docker run
```