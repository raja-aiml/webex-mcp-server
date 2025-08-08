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

// Config holds the application configuration
type Config struct {
	Name        string
	Version     string
	HTTPAddr    string
	EnvPath     string // Path to .env file
	UseAllTools bool   // Whether to load all tools or just core tools
}

// App represents the main application
type App struct {
	config Config
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new application instance
func New(cfg Config) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Run starts the application
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
		if a.config.HTTPAddr != "" {
			errChan <- server.RunHTTPServer(a.ctx, a.config.HTTPAddr, mcpServer, a.config.Name, a.config.Version)
		} else {
			errChan <- server.RunStdioServer(a.ctx, mcpServer, a.config.Name, a.config.Version)
		}
	}()

	// Wait for signal or error
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		return a.Shutdown()
	case err := <-errChan:
		return err
	}
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() error {
	log.Println("Shutting down gracefully...")

	// Cancel context to signal shutdown
	a.cancel()

	// Give some time for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Wait for shutdown or timeout
	select {
	case <-shutdownCtx.Done():
		log.Println("Shutdown timeout exceeded, forcing exit")
	case <-time.After(1 * time.Second):
		log.Println("Shutdown completed")
	}

	return nil
}
