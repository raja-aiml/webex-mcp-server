#!/bin/bash
# Unified MCP server management script - KISS, DRY, YAGNI compliant

set -e

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
IMAGE_NAME="webex-mcp-server:latest"

# Change to project directory
cd "$PROJECT_DIR"

# Load environment
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "Error: .env file not found!"
    echo "Create .env with: WEBEX_PUBLIC_WORKSPACE_API_KEY=your-key"
    exit 1
fi

# Function to build Docker image if needed
build_if_needed() {
    if ! docker images | grep -q "webex-mcp-server.*latest"; then
        echo "Building Docker image..."
        docker build -t "$IMAGE_NAME" .
    fi
}

# Main command handler
case "${1:-help}" in
    stdio)
        build_if_needed
        docker run --rm -i \
            -e WEBEX_PUBLIC_WORKSPACE_API_KEY="$WEBEX_PUBLIC_WORKSPACE_API_KEY" \
            -e WEBEX_API_BASE_URL="${WEBEX_API_BASE_URL:-https://webexapis.com/v1}" \
            "$IMAGE_NAME"
        ;;
    
    inspect)
        build_if_needed
        npx @modelcontextprotocol/inspector docker run --rm -i \
            -e WEBEX_PUBLIC_WORKSPACE_API_KEY="$WEBEX_PUBLIC_WORKSPACE_API_KEY" \
            -e WEBEX_API_BASE_URL="${WEBEX_API_BASE_URL:-https://webexapis.com/v1}" \
            "$IMAGE_NAME"
        ;;
    
    http)
        docker-compose --profile http up --build ${2:--d}
        if [ "$2" = "-d" ] || [ -z "$2" ]; then
            echo "Server running at http://localhost:8084/sse"
            echo "Logs: docker-compose logs -f webex-mcp-http"
            echo "Stop: docker-compose --profile http down"
        fi
        ;;
    
    stop)
        docker-compose --profile http down
        ;;
    
    build)
        docker build -t "$IMAGE_NAME" .
        ;;
    
    help|*)
        cat << EOF
Usage: $(basename $0) [command]

Commands:
  stdio    - Run in stdio mode
  inspect  - Run with MCP Inspector
  http     - Run HTTP/SSE server (add -f for foreground)
  stop     - Stop HTTP server
  build    - Build Docker image
  help     - Show this help

Examples:
  $(basename $0) inspect        # Test with MCP Inspector
  $(basename $0) http          # Start HTTP server in background
  $(basename $0) http -f       # Start HTTP server in foreground
  $(basename $0) stdio         # Run in stdio mode
EOF
        ;;
esac