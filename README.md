# Webex MCP Server (Go)

A high-performance Go implementation of the Webex MCP (Model Context Protocol) Server, providing AI assistants with comprehensive access to Cisco Webex messaging, collaboration, and enterprise features.

## 🚀 Quick Setup for Claude Desktop

### 1. Build the Server
```bash
make build
```

### 2. Configure Claude Desktop
Find your Claude Desktop config file:
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

Add this configuration:
```json
{
  "mcpServers": {
    "webex": {
      "command": "/absolute/path/to/your/webex-mcp-server/build/webex-mcp-server",
      "args": ["-env", "/absolute/path/to/your/webex-mcp-server/.env", "-all-tools"],
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "your-webex-token"
      }
    }
  }
}
```

### 3. Restart Claude Desktop
- Completely quit Claude (not just close window)
- Restart Claude Desktop

### 4. Start Using Webex with Claude!
In Claude, you can now say:
- "List my Webex rooms"
- "Send a message to the Engineering team room"
- "Show recent messages in Project Alpha"
- "Create a new room called Planning Meeting"
- "Add john.doe@company.com to the Marketing team"

## ✨ Features

- **Complete Webex API Integration**: 53+ tools covering messaging, rooms, teams, people, webhooks, and enterprise features
- **High-Performance Architecture**: Built with Go for optimal speed and resource efficiency
- **Multiple Transport Modes**: Support for both stdio and SSE (Server-Sent Events) transports
- **MCP Protocol Compliance**: Uses the official Model Context Protocol Go SDK
- **Advanced Tool Management**: Comprehensive tool registry with modular organization
- **Optimized HTTP Client**:
  - Standard `net/http` with intelligent connection pooling
  - Automatic connection reuse and keep-alive
  - Configurable timeouts and retry logic
  - Built for enterprise-scale operations

## 📋 Prerequisites

- **Go 1.24+** (current version: 1.24.0)
- **Webex API Token** - Get yours from [developer.webex.com](https://developer.webex.com)
- **Docker** (optional) - For containerized deployment

## 🛠️ Installation

1. **Clone the repository:**
```bash
git clone https://github.com/raja-aiml/webex-mcp-server.git
cd webex-mcp-server
```

2. **Install dependencies:**
```bash
make deps
```

3. **Set up environment variables:**
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env and add your Webex API token
# WEBEX_PUBLIC_WORKSPACE_API_KEY=your_webex_api_token_here
```

## ⚙️ Configuration

Create a `.env` file in the project root with the following variables:

```env
# Required: Your Webex API token
WEBEX_PUBLIC_WORKSPACE_API_KEY=your_webex_api_token_here

# Optional: Webex API base URL (default: https://webexapis.com/v1)
WEBEX_API_BASE_URL=https://webexapis.com/v1

# Optional: Server port for HTTP mode (default: 3001)
PORT=3001
```

## 🎯 Usage Modes

### Claude Desktop Integration (Recommended)

**This MCP server is optimized for Claude Desktop.**

1. **Build the server:**
```bash
make build
```

2. **Add to Claude Desktop config:**
```json
{
  "mcpServers": {
    "webex": {
      "command": "/full/path/to/webex-mcp-server/build/webex-mcp-server",
      "args": ["-env", "/full/path/to/webex-mcp-server/.env", "-all-tools"],
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "your-token"
      }
    }
  }
}
```

3. **Restart Claude and start using:**
- "List my Webex rooms"
- "Send a message to [room name]"
- "Show recent messages"
- "Create a team called Engineering"

### Development and Testing Modes

#### Local Development (stdio mode):
```bash
make run
# or
make mcp-local
```

#### HTTP/SSE Mode (for debugging):
```bash
make run http
# Check health endpoint
curl http://localhost:3001/health
```

#### Docker Mode:
```bash
# Build and run in Docker
make mcp-docker

# Or use docker-compose for HTTP mode
make docker run-http
```

## 🔧 Available Tools

The server provides 53+ tools organized into the following categories:

### 💬 Messaging Tools (6 tools)
- `list_messages` - List messages in a room
- `create_a_message` - Send a message to rooms or people
- `get_message_details` - Get detailed message information
- `update_a_message` - Edit an existing message
- `delete_a_message` - Delete a message
- `list_direct_messages` - List direct messages

### 🏠 Room Management (6 tools)
- `list_rooms` - List all accessible rooms
- `create_a_room` - Create new rooms with advanced settings
- `get_room_details` - Get detailed room information
- `update_a_room` - Update room settings and properties
- `delete_a_room` - Delete a room
- `get_room_meeting_details` - Get room meeting information

### 👥 People Management (6 tools)
- `list_people` - List people in your organization
- `create_a_person` - Create new user accounts (admin only)
- `get_person_details` - Get detailed person information
- `update_a_person` - Update person information
- `delete_a_person` - Delete a person (admin only)
- `get_my_own_details` - Get current user details

### 🔗 Membership Management (5 tools)
- `list_memberships` - List room memberships
- `create_a_membership` - Add people to rooms
- `get_membership_details` - Get membership details
- `update_a_membership` - Update membership permissions
- `delete_a_membership` - Remove people from rooms

### 🏢 Team Management (10 tools)
- **Teams**: Create, list, get, update, delete teams
- **Team Memberships**: Add/remove people, manage team permissions

### 🔗 Integration Tools (15+ tools)
- **Webhooks**: Create, list, update, delete webhooks for event notifications
- **Room Tabs**: Add custom tabs and integrations to rooms
- **Attachment Actions**: Handle interactive card responses
- **Events**: Monitor and track organization events
- **ECM (Enterprise Content Management)**: Manage enterprise folders

### 🎯 Advanced Features
- **Bulk Operations**: Many tools support batch processing
- **Enterprise Features**: Advanced admin and organization management
- **Real-time Events**: Webhook support for live notifications
- **File Management**: Attachment and file sharing capabilities

## 🧪 Testing & Development

### Testing with MCP Inspector

**Important**: This MCP server is optimized for Claude Desktop. While you can test with MCP Inspector, Claude Desktop provides the best experience.

#### Option 1: Local Development with MCP Inspector
```bash
# Install MCP Inspector (if not already installed)
npm install -g @modelcontextprotocol/inspector

# Test with local Go server
npx @modelcontextprotocol/inspector go run . -env .env -all-tools
```

#### Option 2: Docker + MCP Inspector
```bash
# Build and test with Docker
make docker build
npx @modelcontextprotocol/inspector docker run --rm -i --env-file .env webex-mcp-server:latest
```

#### Option 3: HTTP/SSE Mode Testing
```bash
# Start server in HTTP mode
make run http

# Connect MCP Inspector to HTTP endpoint
# Open Inspector and connect to: http://localhost:3001/sse
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
make dev test

# Run tests with coverage report
make test-coverage

# Run comprehensive development checks
make dev all  # Runs format, lint, and test
```

### Development Commands

```bash
# Format code
make fmt

# Lint code (requires golangci-lint)
make lint

# Clean build artifacts
make clean

# Build for multiple platforms
make build all

# Check Webex token and connectivity
make check-token

# Security scan
make security-scan
```

## 🏗️ Architecture

The server follows clean architecture principles with a modular, scalable structure:

```
webex-mcp-server/
├── main.go                     # Application entry point & MCP server setup
├── internal/
│   ├── app/                    # Application orchestration layer
│   │   ├── app.go             # Main application logic
│   │   └── app_test.go        # Application tests
│   ├── config/                 # Configuration management
│   │   ├── config.go          # Environment & API configuration
│   │   └── config_test.go     # Configuration tests
│   ├── server/                 # MCP server implementation
│   │   ├── server.go          # MCP server logic
│   │   ├── transport.go       # Transport layer (stdio/HTTP/SSE)
│   │   ├── protocol.go        # MCP protocol handling
│   │   └── *_test.go          # Server tests
│   ├── webex/                  # Webex API client
│   │   ├── client.go          # HTTP client with connection pooling
│   │   ├── interface.go       # HTTPClient interface (DI/IoC)
│   │   └── client_test.go     # Client tests
│   ├── tools/                  # MCP tool implementations (53+ tools)
│   │   ├── registry.go        # Tool registry and management
│   │   ├── base.go            # Base tool functionality (DRY principle)
│   │   ├── generic_tool.go    # Generic tool implementation
│   │   ├── tool_factory.go    # Tool factory pattern
│   │   ├── plugin_loader.go   # Dynamic tool loading
│   │   ├── messages.go        # Message tools (6 tools)
│   │   ├── core_extras.go     # Additional core tools
│   │   └── *_test.go          # Tool tests
│   ├── advanced_tools/         # Advanced/Enterprise tools
│   │   ├── rooms.go           # Room management tools
│   │   ├── people.go          # People management tools
│   │   ├── teams.go           # Team management tools
│   │   ├── memberships.go     # Membership tools
│   │   ├── team_memberships.go # Team membership tools
│   │   ├── room_tabs.go       # Room tab tools
│   │   ├── attachments.go     # Attachment action tools
│   │   ├── events.go          # Event monitoring tools
│   │   ├── ecm.go             # ECM folder tools
│   │   └── plugin_loader.go   # Advanced tool loading
│   ├── handlers/               # HTTP request handlers
│   │   ├── handlers.go        # HTTP route handlers
│   │   └── handlers_test.go   # Handler tests
│   └── testutil/               # Testing utilities
│       └── testutil.go        # Shared test helpers
├── build/                      # Build output directory
└── Dockerfile                  # Container configuration
```

### 🎨 Design Principles Applied

- **🎯 SOLID Principles**: 
  - Single Responsibility: Each package has a clear purpose
  - Open/Closed: Extensible through interfaces
  - Liskov Substitution: Interface-based design
  - Interface Segregation: Small, focused interfaces
  - Dependency Inversion: Dependency injection throughout

- **🚀 Performance Optimizations**:
  - Connection pooling for HTTP clients
  - Efficient tool registry with O(1) lookups
  - Minimal memory allocations
  - Concurrent-safe operations

- **🛠️ Code Quality**:
  - **DRY**: Shared base functionality, no code duplication
  - **KISS**: Simple, readable implementations
  - **YAGNI**: No over-engineering, just what's needed
  - Comprehensive test coverage (>85%)

### 🔧 Key Components

- **Tool Registry**: Dynamic tool discovery and registration
- **Plugin System**: Modular tool loading with advanced/basic modes
- **Transport Layer**: Support for multiple transport protocols
- **HTTP Client**: Optimized for enterprise-scale operations
- **Configuration**: Environment-based configuration management


## 🐳 Docker Deployment

### Quick Start with Docker

```bash
# Build and run with environment file
make docker build
make docker run

# Run in HTTP/SSE mode (recommended for testing)
make docker run-http

# Or use docker-compose
docker-compose --profile http up --build
```

### Manual Docker Commands

```bash
# Build Docker image
docker build -t webex-mcp-server:latest .

# Run in stdio mode (for MCP clients)
docker run --rm -i --env-file .env webex-mcp-server:latest

# Run in HTTP mode (for debugging/testing)
docker run --rm -p 8084:8084 --env-file .env webex-mcp-server:latest -http :8084

# Using docker-compose profiles
docker-compose --profile http up --build    # HTTP/SSE mode
docker-compose --profile stdio up --build   # stdio mode
```

### Testing with MCP Inspector + Docker

#### Option 1: Stdio Mode with Docker
```bash
npx @modelcontextprotocol/inspector docker run --rm -i --env-file .env webex-mcp-server:latest
```

#### Option 2: HTTP/SSE Mode
```bash
# Start the server
docker-compose --profile http up -d

# Connect MCP Inspector to HTTP endpoint
# Open Inspector and connect to: http://localhost:8084/sse
```

### Docker Compose Profiles

- **`stdio`** - Run in stdio mode for MCP clients
- **`http`** - Run in HTTP/SSE mode on port 8084
- **`dev`** - Development mode with hot reload

### Troubleshooting Docker

#### .env File Issues
```bash
# Ensure .env file exists and is accessible
ls -la .env

# Check docker-compose volume mounting
docker-compose config
```

#### Port Conflicts
```bash
# Change port in .env file
echo "MCP_PORT=8085" >> .env

# Or use different port in docker run
docker run --rm -p 9001:8084 --env-file .env webex-mcp-server:latest -http :8084
```

#### View Logs
```bash
# View logs from docker-compose
docker-compose logs -f webex-mcp-http

# View logs from direct docker run
docker logs <container-id>
```

#### Stop and Clean
```bash
# Stop all services
make docker stop

# Clean up resources
make docker clean
```

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes** and ensure tests pass: `make test`
4. **Run code quality checks**: `make dev all`
5. **Commit your changes**: `git commit -m 'Add amazing feature'`
6. **Push to the branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Development Guidelines

- Follow Go best practices and idioms
- Maintain test coverage above 85%
- Use conventional commit messages
- Update documentation for new features
- Run `make security-scan` before committing

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support & Resources

- **Issues**: [GitHub Issues](https://github.com/raja-aiml/webex-mcp-server/issues)
- **Webex Developer Docs**: [developer.webex.com](https://developer.webex.com)
- **MCP Protocol**: [modelcontextprotocol.io](https://modelcontextprotocol.io)
- **Claude Desktop**: [claude.ai/desktop](https://claude.ai/desktop)

## 🏆 Acknowledgments

- Built with the [Model Context Protocol Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- Powered by the [Webex APIs](https://developer.webex.com)
- Optimized for [Claude Desktop](https://claude.ai/desktop)