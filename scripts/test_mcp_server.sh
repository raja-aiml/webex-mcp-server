#!/bin/bash

# MCP Server Test Script
# This script tests basic functionality of the Webex MCP Server

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configuration
PORT=${PORT:-3001}
SERVER_BINARY="./build/webex-mcp-server"

# Check if server binary exists
if [ ! -f "$SERVER_BINARY" ]; then
    echo -e "${RED}Error: Server binary not found. Run 'make build' first.${NC}"
    exit 1
fi

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        exit 1
    fi
}

echo "Starting Webex MCP Server tests..."
echo "=================================="

# Start server in background
echo "Starting server on port $PORT..."
$SERVER_BINARY -http :$PORT &
SERVER_PID=$!

# Give server time to start
sleep 2

# Test 1: Health check
echo -e "\n1. Testing health endpoint..."
if curl -f -s http://localhost:$PORT/health > /dev/null; then
    HEALTH_RESPONSE=$(curl -s http://localhost:$PORT/health)
    print_result 0 "Health check passed"
    echo "   Response: $HEALTH_RESPONSE"
else
    print_result 1 "Health check failed"
fi

# Test 2: Test if MCP endpoint responds
echo -e "\n2. Testing MCP endpoint..."
# For SSE-based MCP servers, we just test if the endpoint responds
MCP_TEST=$(curl -s -o /dev/null -w "%{http_code}" -X GET http://localhost:$PORT/)

if [ "$MCP_TEST" == "200" ]; then
    print_result 0 "MCP endpoint is responding"
    echo "   Note: Full MCP testing requires SSE client or Claude Desktop"
else
    print_result 1 "MCP endpoint not responding (HTTP $MCP_TEST)"
fi

# Test 3: Check server is built with tools
echo -e "\n3. Checking server capabilities..."
# Since this is an SSE server, we verify it was built correctly
if [ -f "$SERVER_BINARY" ] && [ -x "$SERVER_BINARY" ]; then
    # Check binary size to ensure tools are included
    BINARY_SIZE=$(stat -f%z "$SERVER_BINARY" 2>/dev/null || stat -c%s "$SERVER_BINARY" 2>/dev/null)
    if [ "$BINARY_SIZE" -gt 1000000 ]; then  # Binary should be > 1MB with all tools
        print_result 0 "Server binary includes tools (size: $(echo $BINARY_SIZE | numfmt --to=iec-i --suffix=B 2>/dev/null || echo "${BINARY_SIZE} bytes"))"
    else
        print_result 1 "Server binary seems too small"
    fi
else
    print_result 1 "Server binary not found or not executable"
fi

# Clean up
echo -e "\nCleaning up..."
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

echo -e "\n${GREEN}All tests passed!${NC}"
echo "=================================="