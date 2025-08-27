package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raja-aiml/webex-mcp-server/internal/server"
)

// Config holds the application configuration - TypeScript interface compliant
type Config struct {
	Name        string `json:"name"`           // Application name
	Version     string `json:"version"`        // Application version
	HTTPAddr    string `json:"http,omitempty"` // HTTP server address
	EnvPath     string `json:"env,omitempty"`  // Path to .env file
	UseAllTools bool   `json:"useAllTools"`    // Whether to load all tools or just core tools
	SSEMode     bool   `json:"sse,omitempty"`  // Server-Sent Events mode (TypeScript compliance)
}

// App represents the main application - TypeScript class-like structure
type App struct {
	config Config
	ctx    context.Context
	cancel context.CancelFunc
}

// ApplicationState represents the current application state - TypeScript enum-like
type ApplicationState int

const (
	StateInitializing ApplicationState = iota
	StateRunning
	StateShuttingDown
	StateStopped
)

// String method for ApplicationState - TypeScript toString() equivalent
func (s ApplicationState) String() string {
	switch s {
	case StateInitializing:
		return "initializing"
	case StateRunning:
		return "running"
	case StateShuttingDown:
		return "shutting_down"
	case StateStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

// New creates a new application instance - TypeScript constructor pattern
func New(cfg Config) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

// GetConfig returns application configuration - TypeScript getter pattern
func (a *App) GetConfig() Config {
	return a.config
}

// GetState returns current application state - TypeScript getter pattern
func (a *App) GetState() ApplicationState {
	select {
	case <-a.ctx.Done():
		return StateStopped
	default:
		return StateRunning
	}
}

// Run starts the application - TypeScript async method pattern
func (a *App) Run() error {
	// Initialize configuration
	if err := server.InitializeConfig(a.config.EnvPath); err != nil {
		return err
	}

	// Create MCP server with specified mode
	mcpServer, err := server.CreateMCPServerWithMode(a.config.Name, a.config.Version, a.config.UseAllTools)
	if err != nil {
		return err
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		// TypeScript-style conditional execution
		if a.config.HTTPAddr != "" || a.config.SSEMode {
			// Default to port 3001 if SSE mode but no address specified (TypeScript pattern)
			addr := a.config.HTTPAddr
			if addr == "" && a.config.SSEMode {
				addr = ":3001"
			}
			errChan <- server.RunHTTPServer(a.ctx, addr, mcpServer, a.config.Name, a.config.Version)
		} else {
			errChan <- server.RunStdioServer(a.ctx, mcpServer, a.config.Name, a.config.Version)
		}
	}()

	// Wait for signal or error - TypeScript Promise.race() pattern
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		return a.Shutdown()
	case err := <-errChan:
		return err
	}
}

// Shutdown gracefully shuts down the application - TypeScript async method pattern
func (a *App) Shutdown() error {
	log.Println("Shutting down gracefully...")

	// Cancel context to signal shutdown
	a.cancel()

	// Give some time for graceful shutdown - TypeScript timeout pattern
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Wait for shutdown or timeout - TypeScript Promise.race() pattern
	select {
	case <-shutdownCtx.Done():
		log.Println("Shutdown timeout exceeded, forcing exit")
	case <-time.After(1 * time.Second):
		log.Println("Shutdown completed")
	}

	return nil
}
