package main

import (
	"flag"
	"log"

	"github.com/raja-aiml/webex-mcp-server-go/server"
)

const (
	ServerName    = "webex-mcp-server"
	ServerVersion = "0.1.0"
)

func main() {
	var httpAddr string
	flag.StringVar(&httpAddr, "http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
	flag.Parse()

	// Initialize configuration
	if err := server.InitializeConfig(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Create MCP server
	mcpServer, err := server.CreateMCPServer(ServerName, ServerVersion)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Run server in appropriate mode
	if httpAddr != "" {
		if err := server.RunHTTPServer(httpAddr, mcpServer, ServerName, ServerVersion); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	} else {
		if err := server.RunStdioServer(mcpServer, ServerName, ServerVersion); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}
