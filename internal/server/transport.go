package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raja-aiml/webex-mcp-server/internal/handlers"
)

// RunHTTPServer starts the HTTP server with context support
func RunHTTPServer(ctx context.Context, httpAddr string, server *mcp.Server, serviceName, version string) error {
	mux := handlers.SetupHTTPHandlers(server, serviceName, version)

	httpServer := &http.Server{
		Addr:              httpAddr,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		log.Printf("MCP server listening at %s", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	log.Println("Shutting down HTTP server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return httpServer.Shutdown(shutdownCtx)
}

// RunStdioServer starts the server in stdio mode with context support
func RunStdioServer(ctx context.Context, server *mcp.Server, serviceName, version string) error {
	// Create base stdio transport
	var transport mcp.Transport = mcp.NewStdioTransport()

	// Only enable logging if MCP_DEBUG environment variable is set
	if os.Getenv("MCP_DEBUG") == "true" {
		transport = mcp.NewLoggingTransport(transport, os.Stderr)
		log.Println("MCP debug logging enabled (set MCP_DEBUG=false to disable)")
	}

	log.Printf("Starting %s v%s in stdio mode (MCP protocol version: 2024-11-05)", serviceName, version)

	// Run with context
	return server.Run(ctx, transport)
}
