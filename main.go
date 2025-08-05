package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raja-aiml/webex-mcp-server-go/config"
	"github.com/raja-aiml/webex-mcp-server-go/tools"
)

const (
	ServerName    = "webex-mcp-server"
	ServerVersion = "0.1.0"
)

// initializeConfig loads environment variables and validates configuration
func initializeConfig() error {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	return nil
}

// createMCPServer creates and configures the MCP server with tools
func createMCPServer() (*mcp.Server, error) {
	// Load all tools
	toolRegistry, err := tools.LoadTools()
	if err != nil {
		return nil, fmt.Errorf("failed to load tools: %w", err)
	}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    ServerName,
		Version: ServerVersion,
	}, nil)

	// Register all tools with the server
	registerTools(server, toolRegistry)

	return server, nil
}

// healthHandler returns an HTTP handler for health checks
func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"status":  "healthy",
			"service": ServerName,
			"version": ServerVersion,
			"time":    time.Now().UTC().Format(time.RFC3339),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode health response: %v", err)
		}
	}
}

// setupHTTPHandlers configures HTTP handlers for the server
func setupHTTPHandlers(server *mcp.Server) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", healthHandler())

	// MCP handler
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)
	mux.Handle("/", handler)

	return mux
}

// runHTTPServer starts the HTTP server
func runHTTPServer(httpAddr string, server *mcp.Server) error {
	mux := setupHTTPHandlers(server)
	log.Printf("MCP server listening at %s (using fasthttp client for Webex API)", httpAddr)
	return http.ListenAndServe(httpAddr, mux)
}

// runStdioServer starts the server in stdio mode
func runStdioServer(server *mcp.Server) error {
	transport := mcp.NewLoggingTransport(mcp.NewStdioTransport(), os.Stderr)
	log.Printf("Starting %s v%s in stdio mode", ServerName, ServerVersion)
	return server.Run(context.Background(), transport)
}

func main() {
	var httpAddr string
	flag.StringVar(&httpAddr, "http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
	flag.Parse()

	// Initialize configuration
	if err := initializeConfig(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Create MCP server
	server, err := createMCPServer()
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Run server in appropriate mode
	if httpAddr != "" {
		if err := runHTTPServer(httpAddr, server); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	} else {
		if err := runStdioServer(server); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}

// convertToolSchema converts various schema formats to jsonschema.Schema
func convertToolSchema(tool tools.Tool) (*jsonschema.Schema, error) {
	schemaInterface := tool.GetInputSchema()

	// Handle both direct jsonschema.Schema and legacy interface{} schemas
	switch s := schemaInterface.(type) {
	case *jsonschema.Schema:
		return s, nil
	default:
		// Legacy path for tools still using interface{} schemas
		schemaJSON, err := json.Marshal(schemaInterface)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal schema: %w", err)
		}
		schema := &jsonschema.Schema{}
		if err := json.Unmarshal(schemaJSON, schema); err != nil {
			return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
		}
		return schema, nil
	}
}

// createToolHandler creates an MCP tool handler for a given tool
func createToolHandler(tool tools.Tool) mcp.ToolHandler {
	return func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
		// Convert arguments to JSON for our tool
		argsJSON, err := json.Marshal(params.Arguments)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal arguments: %w", err)
		}

		// Execute the tool with raw arguments
		result, err := tool.Execute(argsJSON)
		if err != nil {
			// Return error as tool result per MCP spec
			// Tool errors should be in the result, not as protocol errors
			return &mcp.CallToolResultFor[any]{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error: %v", err),
					},
				},
				IsError: true,
			}, nil
		}

		// Convert result to JSON string
		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}

		// Return as text content
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(resultJSON),
				},
			},
		}, nil
	}
}

// registerTools registers all tools from the registry with the MCP server
func registerTools(server *mcp.Server, registry *tools.Registry) {
	allTools := registry.GetTools()

	for _, tool := range allTools {
		// Convert schema
		schema, err := convertToolSchema(tool)
		if err != nil {
			log.Printf("Failed to convert schema for tool %s: %v", tool.Name(), err)
			continue
		}

		// Create MCP tool definition
		mcpTool := &mcp.Tool{
			Name:        tool.Name(),
			Description: tool.Description(),
			InputSchema: schema,
		}

		// Create and add handler
		handler := createToolHandler(tool)
		server.AddTool(mcpTool, handler)
	}
}
