package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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

func main() {
	var httpAddr string
	flag.StringVar(&httpAddr, "http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Load all tools
	toolRegistry, err := tools.LoadTools()
	if err != nil {
		log.Fatalf("Failed to load tools: %v", err)
	}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    ServerName,
		Version: ServerVersion,
	}, nil)

	// Register all tools with the server
	registerTools(server, toolRegistry)

	if httpAddr != "" {
		// HTTP/SSE mode
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		log.Printf("MCP server listening at %s (using fasthttp client for Webex API)", httpAddr)
		if err := http.ListenAndServe(httpAddr, handler); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	} else {
		// Stdio mode
		transport := mcp.NewLoggingTransport(mcp.NewStdioTransport(), os.Stderr)
		log.Printf("Starting %s v%s in stdio mode", ServerName, ServerVersion)
		if err := server.Run(context.Background(), transport); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}

func registerTools(server *mcp.Server, registry *tools.Registry) {
	// Get all tools from registry
	allTools := registry.GetTools()
	
	for _, tool := range allTools {
		// Capture tool in closure
		t := tool
		
		// Get schema directly - simplified approach
		schemaInterface := t.GetInputSchema()
		
		// Handle both direct jsonschema.Schema and legacy interface{} schemas
		var schema *jsonschema.Schema
		switch s := schemaInterface.(type) {
		case *jsonschema.Schema:
			schema = s
		default:
			// Legacy path for tools still using interface{} schemas
			schemaJSON, err := json.Marshal(schemaInterface)
			if err != nil {
				log.Printf("Failed to marshal schema for tool %s: %v", t.Name(), err)
				continue
			}
			schema = &jsonschema.Schema{}
			if err := json.Unmarshal(schemaJSON, schema); err != nil {
				log.Printf("Failed to unmarshal schema for tool %s: %v", t.Name(), err)
				continue
			}
		}
		
		// Create MCP tool definition
		mcpTool := &mcp.Tool{
			Name:        t.Name(),
			Description: t.Description(),
			InputSchema: schema,
		}
		
		// Create handler function using the proper ToolHandler type
		handler := mcp.ToolHandler(func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
			// Convert arguments to JSON for our tool
			argsJSON, err := json.Marshal(params.Arguments)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal arguments: %w", err)
			}
			
			// Execute the tool with raw arguments
			result, err := t.Execute(argsJSON)
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
		})
		
		// Add tool to server
		server.AddTool(mcpTool, handler)
	}
}