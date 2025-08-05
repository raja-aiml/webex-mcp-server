#!/usr/bin/env python3
"""
Test script for MCP Server using SSE (Server-Sent Events)
"""

import json
import requests
import sys
from typing import Dict, Any

# Configuration
SERVER_URL = "http://localhost:3001"

def test_mcp_request(method: str, params: Dict[str, Any] = None) -> Dict[str, Any]:
    """Send an MCP request and get response"""
    request = {
        "jsonrpc": "2.0",
        "method": method,
        "id": 1
    }
    if params:
        request["params"] = params
    
    # For SSE-based servers, we need to handle streaming
    headers = {
        "Content-Type": "application/json",
        "Accept": "text/event-stream"
    }
    
    try:
        # Send request
        response = requests.post(
            SERVER_URL,
            json=request,
            headers=headers,
            stream=True,
            timeout=5
        )
        
        # Parse SSE response
        result = None
        for line in response.iter_lines():
            if line:
                line = line.decode('utf-8')
                if line.startswith('data: '):
                    data = line[6:]  # Remove 'data: ' prefix
                    try:
                        result = json.loads(data)
                        if 'id' in result and result['id'] == 1:
                            return result
                    except json.JSONDecodeError:
                        continue
        
        return result or {"error": "No valid response received"}
        
    except requests.exceptions.RequestException as e:
        return {"error": f"Request failed: {str(e)}"}

def main():
    print("=== MCP Server SSE Test ===\n")
    
    # Test 1: Check health
    print("1. Testing health endpoint...")
    try:
        response = requests.get(f"{SERVER_URL}/health")
        if response.status_code == 200:
            print(f"✓ Health check passed: {response.json()}")
        else:
            print(f"✗ Health check failed: {response.status_code}")
    except Exception as e:
        print(f"✗ Health check error: {e}")
    
    # Test 2: List tools
    print("\n2. Testing tools/list...")
    result = test_mcp_request("tools/list")
    if "result" in result:
        tools = result["result"]
        print(f"✓ Found {len(tools)} tools")
        # Show first 5 tools
        for i, tool in enumerate(tools[:5]):
            print(f"   - {tool.get('name', 'unknown')}")
        if len(tools) > 5:
            print(f"   ... and {len(tools) - 5} more")
    else:
        print(f"✗ Failed to list tools: {result}")
    
    # Test 3: Call a tool
    print("\n3. Testing tool call (get_my_own_details)...")
    result = test_mcp_request("tools/call", {
        "name": "get_my_own_details",
        "arguments": {}
    })
    if "result" in result:
        print("✓ Tool call successful")
        print(f"   Response: {json.dumps(result['result'], indent=2)[:200]}...")
    elif "error" in result:
        # This is expected if no API key is set
        error = result.get("error", {})
        if "message" in error and "401" in str(error["message"]):
            print("✓ Tool call returned expected auth error (no API key)")
        else:
            print(f"✗ Tool call error: {error}")
    else:
        print(f"✗ Unexpected response: {result}")
    
    print("\n=== Test Complete ===")

if __name__ == "__main__":
    main()