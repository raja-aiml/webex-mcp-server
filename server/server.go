package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raja-aiml/webex-mcp-server-go/tools"
)

// CreateMCPServer creates and configures the MCP server with tools
func CreateMCPServer(name, version string) (*mcp.Server, error) {
	// Load all tools
	toolRegistry, err := tools.LoadTools()
	if err != nil {
		return nil, fmt.Errorf("failed to load tools: %w", err)
	}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    name,
		Version: version,
	}, nil)

	// Register all tools with the server
	registerTools(server, toolRegistry)

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
