#!/bin/bash

# Interactive MCP Client for testing the Webex MCP Server
# This script provides an easy way to test MCP server commands

# Configuration
SERVER_URL="${SERVER_URL:-http://localhost:3001}"
REQUEST_ID=1

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Function to send MCP request
send_request() {
    local method=$1
    local params=$2
    local id=$REQUEST_ID
    
    if [ -z "$params" ]; then
        local request="{\"jsonrpc\": \"2.0\", \"method\": \"$method\", \"id\": $id}"
    else
        local request="{\"jsonrpc\": \"2.0\", \"method\": \"$method\", \"params\": $params, \"id\": $id}"
    fi
    
    echo -e "${BLUE}Request:${NC}"
    echo "$request" | jq .
    
    echo -e "\n${GREEN}Response:${NC}"
    curl -s -X POST "$SERVER_URL" \
        -H "Content-Type: application/json" \
        -d "$request" | jq .
    
    ((REQUEST_ID++))
}

# Function to call a tool
call_tool() {
    local tool_name=$1
    local args=$2
    
    local params="{\"name\": \"$tool_name\", \"arguments\": $args}"
    send_request "tools/call" "$params"
}

# Print header
clear
echo -e "${YELLOW}=== Webex MCP Server Interactive Client ===${NC}"
echo -e "Server: $SERVER_URL\n"

# Check if server is running
echo "Checking server health..."
if curl -f -s "$SERVER_URL/health" > /dev/null; then
    echo -e "${GREEN}✓ Server is healthy${NC}\n"
else
    echo -e "${RED}✗ Server is not responding. Start it with: make run http${NC}"
    exit 1
fi

# Main menu
while true; do
    echo -e "\n${BLUE}Choose an option:${NC}"
    echo "1. List all tools"
    echo "2. List rooms"
    echo "3. List people"
    echo "4. Get my details"
    echo "5. Send a message"
    echo "6. Custom tool call"
    echo "7. Raw MCP request"
    echo "0. Exit"
    
    read -p "Enter choice: " choice
    
    case $choice in
        1)
            echo -e "\n${YELLOW}Listing all tools...${NC}"
            send_request "tools/list"
            ;;
        2)
            echo -e "\n${YELLOW}Listing rooms...${NC}"
            call_tool "list_rooms" '{"max": 10}'
            ;;
        3)
            echo -e "\n${YELLOW}Listing people...${NC}"
            call_tool "list_people" '{"max": 10}'
            ;;
        4)
            echo -e "\n${YELLOW}Getting my details...${NC}"
            call_tool "get_my_own_details" '{}'
            ;;
        5)
            read -p "Enter room ID: " room_id
            read -p "Enter message text: " message
            echo -e "\n${YELLOW}Sending message...${NC}"
            call_tool "create_a_message" "{\"roomId\": \"$room_id\", \"text\": \"$message\"}"
            ;;
        6)
            read -p "Enter tool name: " tool_name
            echo "Enter arguments as JSON (or press Enter for empty):"
            read args
            if [ -z "$args" ]; then
                args="{}"
            fi
            echo -e "\n${YELLOW}Calling tool: $tool_name${NC}"
            call_tool "$tool_name" "$args"
            ;;
        7)
            read -p "Enter MCP method: " method
            echo "Enter params as JSON (or press Enter for none):"
            read params
            echo -e "\n${YELLOW}Sending raw request...${NC}"
            send_request "$method" "$params"
            ;;
        0)
            echo "Goodbye!"
            exit 0
            ;;
        *)
            echo "Invalid choice. Please try again."
            ;;
    esac
    
    echo -e "\n${YELLOW}Press Enter to continue...${NC}"
    read
done