#!/bin/bash
# Start the MCP server in HTTP/SSE mode

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    echo "Please create a .env file with your WEBEX_PUBLIC_WORKSPACE_API_KEY"
    exit 1
fi

# Source environment variables
source .env

# Build and start the HTTP server
echo "Building and starting MCP server in HTTP mode..."
docker-compose --profile http up --build -d

# Wait for server to be ready
echo "Waiting for server to be ready..."
sleep 3

# Check if server is running
if docker-compose ps | grep -q "webex-mcp-http.*Up"; then
    echo "✅ MCP server is running in HTTP mode at http://localhost:8084"
    echo ""
    echo "To test with MCP Inspector:"
    echo "  1. Open MCP Inspector"
    echo "  2. Connect to: http://localhost:8084/sse"
    echo ""
    echo "To view logs:"
    echo "  docker-compose logs -f webex-mcp-http"
    echo ""
    echo "To stop the server:"
    echo "  docker-compose --profile http down"
else
    echo "❌ Failed to start MCP server"
    echo "Check logs with: docker-compose logs webex-mcp-http"
    exit 1
fi