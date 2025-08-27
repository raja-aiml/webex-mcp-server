package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

// CreateMCPServer creates and configures the MCP server with tools
// By default, loads only core tools for minimal conversation functionality
func CreateMCPServer(name, version string) (*mcp.Server, error) {
	return CreateMCPServerWithMode(name, version, false)
}

// CreateMCPServerWithMode creates MCP server with specified tool mode
// useAllTools: false = core tools only (minimal), true = all tools (full functionality)
func CreateMCPServerWithMode(name, version string, useAllTools bool) (*mcp.Server, error) {
	var toolRegistry *tools.Registry
	var err error

	// Load tools based on mode
	if useAllTools {
		log.Println("Loading all tools (core + advanced)")
		toolRegistry, err = tools.LoadAllTools()
	} else {
		log.Println("Loading core tools only (minimal mode)")
		toolRegistry, err = tools.LoadCoreTools()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load tools: %w", err)
	}

	// Create MCP server with proper options
	server := mcp.NewServer(&mcp.Implementation{
		Name:    name,
		Version: version,
	}, &mcp.ServerOptions{
		Instructions: fmt.Sprintf("%s v%s - A Model Context Protocol server for Webex messaging operations", name, version),
	})

	// Register all tools with the server
	registerTools(server, toolRegistry)

	// Log loaded tools count
	log.Printf("Loaded %d tools", len(toolRegistry.GetTools()))

	return server, nil
}

// convertToolSchema converts various schema formats to jsonschema.Schema
func convertToolSchema(tool tools.Tool) (*jsonschema.Schema, error) {
	schemaInterface := tool.GetInputSchema()

	// Handle both direct jsonschema.Schema and legacy any schemas
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
	return func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Validate required arguments at handler level
		arguments := request.Params.Arguments
		if arguments == nil {
			arguments = make(map[string]any)
		}

		// Convert arguments to JSON for our tool
		argsJSON, err := json.Marshal(arguments)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Invalid arguments: %v", err),
					},
				},
				IsError: true,
			}, nil
		}

		// Execute the tool with raw arguments
		result, err := tool.Execute(argsJSON)
		if err != nil {
			// Return error as tool result per MCP spec
			// Be specific about the error type
			errorMessage := fmt.Sprintf("Tool execution failed: %v", err)
			if err.Error() == "404 Not Found" {
				errorMessage = "Resource not found. Please verify the ID or name."
			} else if err.Error() == "401 Unauthorized" {
				errorMessage = "Authentication failed. Please check your API credentials."
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: errorMessage,
					},
				},
				IsError: true,
			}, nil
		}

		// Handle different result types
		var content []mcp.Content

		switch v := result.(type) {
		case nil:
			// Empty result (e.g., from DELETE operations)
			content = []mcp.Content{
				&mcp.TextContent{
					Text: "Operation completed successfully",
				},
			}
		case string:
			content = []mcp.Content{
				&mcp.TextContent{
					Text: v,
				},
			}
		case map[string]interface{}:
			// Check if it's a simple success response
			if success, ok := v["success"].(bool); ok && success && len(v) == 1 {
				content = []mcp.Content{
					&mcp.TextContent{
						Text: "Operation completed successfully",
					},
				}
			} else {
				// Format as JSON for complex objects
				resultJSON, err := json.MarshalIndent(v, "", "  ")
				if err != nil {
					return &mcp.CallToolResult{
						Content: []mcp.Content{
							&mcp.TextContent{
								Text: fmt.Sprintf("Failed to format result: %v", err),
							},
						},
						IsError: true,
					}, nil
				}
				content = []mcp.Content{
					&mcp.TextContent{
						Text: string(resultJSON),
					},
				}
			}
		default:
			// Default to JSON serialization
			resultJSON, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf("Failed to format result: %v", err),
						},
					},
					IsError: true,
				}, nil
			}
			content = []mcp.Content{
				&mcp.TextContent{
					Text: string(resultJSON),
				},
			}
		}

		// Return successful result
		return &mcp.CallToolResult{
			Content: content,
			IsError: false,
		}, nil
	}
}

// registerTools registers all tools from the registry with the MCP server
func registerTools(server *mcp.Server, registry *tools.Registry) {
	allTools := registry.GetTools()

	for _, tool := range allTools {
		// Validate tool name for MCP compliance
		if err := ValidateToolName(tool.Name()); err != nil {
			log.Printf("Skipping tool %s: %v", tool.Name(), err)
			continue
		}

		// Validate tool description for MCP compliance
		if err := ValidateToolDescription(tool.Description()); err != nil {
			log.Printf("Skipping tool %s: %v", tool.Name(), err)
			continue
		}

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
