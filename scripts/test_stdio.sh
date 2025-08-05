#!/bin/bash

# Simple stdio test for MCP server
# This shows how the server actually communicates

echo "Testing Webex MCP Server stdio mode..."
echo "======================================"

# Create a temporary file for the request
REQUEST_FILE=$(mktemp)
RESPONSE_FILE=$(mktemp)

# Initialize request (MCP protocol requires this first)
cat > "$REQUEST_FILE" << 'EOF'
{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"tools": {}}}, "id": 1}
EOF

echo "1. Sending initialize request..."
./build/webex-mcp-server < "$REQUEST_FILE" 2>/dev/null | head -1 > "$RESPONSE_FILE"

if grep -q '"result"' "$RESPONSE_FILE"; then
    echo "✓ Server initialized successfully"
    cat "$RESPONSE_FILE" | jq .
else
    echo "✗ Failed to initialize server"
    cat "$RESPONSE_FILE"
fi

# Clean up
rm -f "$REQUEST_FILE" "$RESPONSE_FILE"

echo ""
echo "Note: The server uses stdio for MCP protocol communication."
echo "For full testing, use Claude Desktop or implement a proper MCP client."