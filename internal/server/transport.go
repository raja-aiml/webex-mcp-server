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
	errChan := make(chan error, 1)
	go func() {
		log.Printf("[%s v%s] MCP server listening at %s", serviceName, version, httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for context or error
	select {
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutdownCtx)
	case err := <-errChan:
		return err
	}
}

// RunStdioServer starts the server in stdio mode with context support
func RunStdioServer(ctx context.Context, server *mcp.Server, serviceName, version string) error {
	var transport mcp.Transport = mcp.NewStdioTransport()

	if os.Getenv("MCP_DEBUG") == "true" {
		transport = mcp.NewLoggingTransport(transport, os.Stderr)
		log.Printf("[%s v%s] MCP debug logging enabled", serviceName, version)
	}

	log.Printf("[%s v%s] Starting in stdio mode (MCP protocol version: %s)", serviceName, version, MCPProtocolVersion)
	return server.Run(ctx, transport)
}

// RunSSEServer starts the server in SSE mode
func RunSSEServer(ctx context.Context, httpAddr string, server *mcp.Server, serviceName, version string) error {
	if httpAddr == "" {
		httpAddr = ":3001"
	}

	log.Printf("[%s v%s] Starting SSE server at %s", serviceName, version, httpAddr)
	return RunHTTPServer(ctx, httpAddr, server, serviceName, version)
}
