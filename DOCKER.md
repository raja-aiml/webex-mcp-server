# Docker Setup for Webex MCP Server

## Quick Start

### Using Helper Scripts

```bash
# Test with MCP Inspector (stdio mode)
./inspect-docker.sh

# Start HTTP/SSE server
./start-http-server.sh

# Run with docker-compose
./run-mcp-http.sh
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