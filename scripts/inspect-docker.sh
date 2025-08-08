#!/bin/bash
# Script to run MCP server with Docker for inspection

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    echo "Please create a .env file with your WEBEX_PUBLIC_WORKSPACE_API_KEY"
    exit 1
fi

# Source environment variables
source .env

# Build the Docker image if needed
echo "Building Docker image..."
docker build -t webex-mcp-server:latest .

# Run the MCP Inspector with Docker
echo "Starting MCP Inspector with Docker..."
npx @modelcontextprotocol/inspector docker run --rm -i --env-file .env webex-mcp-server:latest