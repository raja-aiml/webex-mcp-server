package server

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raja-aiml/webex-mcp-server-go/handlers"
)

// RunHTTPServer starts the HTTP server
func RunHTTPServer(httpAddr string, server *mcp.Server, serviceName, version string) error {
	mux := handlers.SetupHTTPHandlers(server, serviceName, version)
	log.Printf("MCP server listening at %s (using fasthttp client for Webex API)", httpAddr)
	return http.ListenAndServe(httpAddr, mux)
}

// RunStdioServer starts the server in stdio mode
func RunStdioServer(server *mcp.Server, serviceName, version string) error {
	transport := mcp.NewLoggingTransport(mcp.NewStdioTransport(), os.Stderr)
	log.Printf("Starting %s v%s in stdio mode", serviceName, version)
	return server.Run(context.Background(), transport)
}
