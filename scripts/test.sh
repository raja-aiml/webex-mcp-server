#!/bin/bash
# Unified test script for Webex MCP Server

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Usage
usage() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  health    - Test health endpoint (HTTP mode)"
    echo "  stdio     - Test stdio initialization"
    echo "  bench     - Run performance benchmark"
    echo "  all       - Run all tests (default)"
    echo ""
    exit 1
}

# Print result
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
    fi
}

# Test health endpoint
test_health() {
    echo -e "${YELLOW}Testing health endpoint...${NC}"
    
    # Start server
    ./build/webex-mcp-server -http :3001 &
    SERVER_PID=$!
    sleep 2
    
    # Test health
    if curl -f -s http://localhost:3001/health > /dev/null; then
        RESPONSE=$(curl -s http://localhost:3001/health)
        print_result 0 "Health check passed"
        echo "  Response: $RESPONSE"
    else
        print_result 1 "Health check failed"
    fi
    
    # Cleanup
    kill $SERVER_PID 2>/dev/null || true
}

# Test stdio mode
test_stdio() {
    echo -e "${YELLOW}Testing stdio mode...${NC}"
    
    # Test initialization
    RESPONSE=$(echo '{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}}, "id": 1}' | ./build/webex-mcp-server 2>/dev/null | head -1)
    
    if echo "$RESPONSE" | grep -q '"result"'; then
        print_result 0 "Server initialized successfully"
        echo "$RESPONSE" | jq -r '"\n  Protocol: \(.result.protocolVersion)\n  Server: \(.result.serverInfo.name) v\(.result.serverInfo.version)\n  Tools: \(.result.capabilities.tools.listChanged)"' 2>/dev/null || true
    else
        print_result 1 "Failed to initialize server"
    fi
}

# Run benchmark
test_bench() {
    echo -e "${YELLOW}Running benchmark...${NC}"
    
    # Start server
    ./build/webex-mcp-server -http :3001 &
    SERVER_PID=$!
    sleep 2
    
    # Check if ab is available
    if command -v ab &> /dev/null; then
        echo "Running 1000 requests with concurrency 10..."
        ab -n 1000 -c 10 -q http://localhost:3001/health 2>&1 | grep -E "(Requests per second:|Time per request:|Transfer rate:)"
    else
        echo "Apache Bench (ab) not installed. Skipping benchmark."
    fi
    
    # Cleanup
    kill $SERVER_PID 2>/dev/null || true
}

# Main
main() {
    # Check if binary exists
    if [ ! -f "./build/webex-mcp-server" ]; then
        echo -e "${RED}Error: Server binary not found. Run 'make build' first.${NC}"
        exit 1
    fi
    
    # Parse command
    case ${1:-all} in
        health)
            test_health
            ;;
        stdio)
            test_stdio
            ;;
        bench)
            test_bench
            ;;
        all)
            test_stdio
            echo ""
            test_health
            echo ""
            test_bench
            ;;
        *)
            usage
            ;;
    esac
    
    echo -e "\n${GREEN}Testing complete!${NC}"
}

main "$@"