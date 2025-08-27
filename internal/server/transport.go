package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raja-aiml/webex-mcp-server/internal/handlers"
)

// TransportType represents different transport types - TypeScript enum pattern
type TransportType string

const (
	TransportStdio TransportType = "stdio"
	TransportHTTP  TransportType = "http"
	TransportSSE   TransportType = "sse"
)

// String method for TransportType - TypeScript toString() equivalent
func (t TransportType) String() string {
	return string(t)
}

// TransportConfig represents transport configuration - TypeScript interface pattern
type TransportConfig struct {
	Type    TransportType `json:"type"`
	Address string        `json:"address,omitempty"`
	Port    string        `json:"port,omitempty"`
	Timeout time.Duration `json:"timeout,omitempty"`
}

// ServerOptions represents server options - TypeScript options pattern
type ServerOptions struct {
	Name    string          `json:"name"`
	Version string          `json:"version"`
	Config  TransportConfig `json:"config"`
}

// HTTPServerConfig represents HTTP server configuration - TypeScript config pattern
type HTTPServerConfig struct {
	ReadTimeout       time.Duration `json:"readTimeout"`
	ReadHeaderTimeout time.Duration `json:"readHeaderTimeout"`
	WriteTimeout      time.Duration `json:"writeTimeout"`
	IdleTimeout       time.Duration `json:"idleTimeout"`
	MaxHeaderBytes    int           `json:"maxHeaderBytes"`
}

// getDefaultHTTPConfig returns default HTTP configuration - TypeScript default config pattern
func getDefaultHTTPConfig() HTTPServerConfig {
	return HTTPServerConfig{
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}
}

// RunHTTPServer starts the HTTP server with context support - TypeScript async function pattern
func RunHTTPServer(ctx context.Context, httpAddr string, server *mcp.Server, serviceName, version string) error {
	// TypeScript-style configuration with defaults
	config := getDefaultHTTPConfig()

	// Create HTTP handler using TypeScript-style dependency injection
	mux := handlers.SetupHTTPHandlers(server, serviceName, version)

	// Create HTTP server with TypeScript-style options pattern
	httpServer := &http.Server{
		Addr:              httpAddr,
		Handler:           mux,
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
		MaxHeaderBytes:    config.MaxHeaderBytes,
	}

	// Start server in goroutine - TypeScript Promise pattern
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("[%s v%s] MCP server listening at %s", serviceName, version, httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// TypeScript-style Promise.race() pattern - wait for context or server error
	select {
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down HTTP server...")
	case err := <-serverErrors:
		log.Printf("HTTP server error: %v", err)
		return err
	}

	// Graceful shutdown with timeout - TypeScript async/await pattern
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return httpServer.Shutdown(shutdownCtx)
}

// RunStdioServer starts the server in stdio mode with context support - TypeScript async function pattern
func RunStdioServer(ctx context.Context, server *mcp.Server, serviceName, version string) error {
	// Create base stdio transport - TypeScript factory pattern
	var transport mcp.Transport = mcp.NewStdioTransport()

	// TypeScript-style environment variable handling
	debugMode := getEnvBool("MCP_DEBUG", false)

	// Only enable logging if MCP_DEBUG environment variable is set - TypeScript conditional pattern
	if debugMode {
		transport = mcp.NewLoggingTransport(transport, os.Stderr)
		log.Printf("[%s v%s] MCP debug logging enabled (set MCP_DEBUG=false to disable)", serviceName, version)
	}

	// TypeScript-style logging with template literals equivalent
	log.Printf("[%s v%s] Starting in stdio mode (MCP protocol version: 2024-11-05)", serviceName, version)

	// Run with context - TypeScript async/await pattern
	return server.Run(ctx, transport)
}

// RunSSEServer starts the server in SSE mode - TypeScript async function pattern
func RunSSEServer(ctx context.Context, httpAddr string, server *mcp.Server, serviceName, version string) error {
	// TypeScript-style default parameter handling
	if httpAddr == "" {
		httpAddr = ":3001"
	}

	log.Printf("[%s v%s] Starting SSE server at %s", serviceName, version, httpAddr)

	// Use standard HTTP server for SSE - TypeScript reuse pattern
	return RunHTTPServer(ctx, httpAddr, server, serviceName, version)
}

// getEnvBool gets environment variable as boolean - TypeScript utility function pattern
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// TypeScript-style boolean conversion
	switch value {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}

// TransportManager manages different transport types - TypeScript class pattern
type TransportManager struct {
	options ServerOptions
	server  *mcp.Server
	ctx     context.Context
}

// NewTransportManager creates a new transport manager - TypeScript constructor pattern
func NewTransportManager(options ServerOptions, server *mcp.Server, ctx context.Context) *TransportManager {
	return &TransportManager{
		options: options,
		server:  server,
		ctx:     ctx,
	}
}

// Start starts the appropriate transport - TypeScript async method pattern
func (tm *TransportManager) Start() error {
	switch tm.options.Config.Type {
	case TransportHTTP:
		return RunHTTPServer(tm.ctx, tm.options.Config.Address, tm.server, tm.options.Name, tm.options.Version)
	case TransportSSE:
		return RunSSEServer(tm.ctx, tm.options.Config.Address, tm.server, tm.options.Name, tm.options.Version)
	case TransportStdio:
		return RunStdioServer(tm.ctx, tm.server, tm.options.Name, tm.options.Version)
	default:
		return &TransportError{
			Type:      "unsupported",
			Message:   "unsupported transport type",
			Transport: string(tm.options.Config.Type),
		}
	}
}

// TransportError represents transport-related errors - TypeScript error pattern
type TransportError struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Transport string `json:"transport,omitempty"`
}

// Error implements the error interface - TypeScript Error pattern
func (e *TransportError) Error() string {
	if e.Transport != "" {
		return fmt.Sprintf("%s transport error: %s (transport: %s)", e.Type, e.Message, e.Transport)
	}
	return fmt.Sprintf("%s transport error: %s", e.Type, e.Message)
}
