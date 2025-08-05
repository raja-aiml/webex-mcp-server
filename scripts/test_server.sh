#!/bin/bash

# Test script for Webex MCP Server

echo "Testing Webex MCP Server..."
echo

# Test 1: Initialize
echo "Test 1: Initialize"
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"0.1.0","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | ./build/webex-mcp-server 2>&1 | jq '.'
echo

# Test 2: List tools
echo "Test 2: List tools"
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./build/webex-mcp-server 2>&1 | jq '.result.tools | length'
echo

# Test 3: Call get_my_own_details tool (if API key is set)
if [ -f .env ]; then
    source .env
    if [ ! -z "$WEBEX_PUBLIC_WORKSPACE_API_KEY" ]; then
        echo "Test 3: Call get_my_own_details tool"
        echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_my_own_details","arguments":{}}}' | ./build/webex-mcp-server 2>&1 | jq '.'
    else
        echo "Test 3: Skipped (no API key)"
    fi
else
    echo "Test 3: Skipped (no .env file)"
fi