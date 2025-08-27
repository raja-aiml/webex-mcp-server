package server

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raja-aiml/webex-mcp-server/internal/testutil"
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

func TestCreateMCPServer(t *testing.T) {
	// Set up required environment variable
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-key")
	defer cleanup()

	tests := []struct {
		name    string
		srvName string
		version string
		wantErr bool
	}{
		{
			name:    "creates server successfully",
			srvName: "test-server",
			version: "1.0.0",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := CreateMCPServer(tt.srvName, tt.version)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMCPServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && server == nil {
				t.Error("CreateMCPServer() returned nil server")
			}
		})
	}
}

// mockTool implements the Tool interface for testing
type mockTool struct {
	name        string
	description string
	schemaType  string // "jsonschema" or "legacy"
	executeErr  error
	executeResp interface{}
}

func (m *mockTool) Name() string {
	return m.name
}

func (m *mockTool) Description() string {
	return m.description
}

func (m *mockTool) GetInputSchema() interface{} {
	if m.schemaType == "jsonschema" {
		return &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"test": {Type: "string"},
			},
		}
	}
	// Legacy schema format
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"test": map[string]interface{}{"type": "string"},
		},
	}
}

func (m *mockTool) Execute(input json.RawMessage) (interface{}, error) {
	if m.executeErr != nil {
		return nil, m.executeErr
	}
	return m.executeResp, nil
}

func (m *mockTool) ExecuteWithMap(input map[string]interface{}) (interface{}, error) {
	return m.Execute(nil)
}

func TestConvertToolSchema(t *testing.T) {
	tests := []struct {
		name       string
		tool       tools.Tool
		wantErr    bool
		checkProps bool
	}{
		{
			name: "converts jsonschema.Schema directly",
			tool: &mockTool{
				schemaType: "jsonschema",
			},
			wantErr:    false,
			checkProps: true,
		},
		{
			name: "converts legacy schema format",
			tool: &mockTool{
				schemaType: "legacy",
			},
			wantErr:    false,
			checkProps: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, err := convertToolSchema(tt.tool)

			if (err != nil) != tt.wantErr {
				t.Errorf("convertToolSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && schema == nil {
				t.Error("convertToolSchema() returned nil schema")
			}

			if tt.checkProps && !tt.wantErr {
				if schema.Type != "object" {
					t.Errorf("Schema type = %v, want object", schema.Type)
				}
				if _, ok := schema.Properties["test"]; !ok {
					t.Error("Schema missing expected 'test' property")
				}
			}
		})
	}
}

func TestCreateToolHandler(t *testing.T) {
	tests := []struct {
		name      string
		tool      *mockTool
		args      map[string]any
		wantError bool
		toolError bool
	}{
		{
			name: "successful execution",
			tool: &mockTool{
				name:        "test-tool",
				executeResp: map[string]string{"result": "success"},
			},
			args:      map[string]any{"input": "test"},
			wantError: false,
			toolError: false,
		},
		{
			name: "tool execution error",
			tool: &mockTool{
				name:       "test-tool",
				executeErr: fmt.Errorf("tool error"),
			},
			args:      map[string]any{"input": "test"},
			wantError: false,
			toolError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := createToolHandler(tt.tool)

			// Create a mock server session
			_ = mcp.NewServer(&mcp.Implementation{
				Name:    "test",
				Version: "1.0.0",
			}, nil)

			// Create request
			request := &mcp.CallToolRequest{
				Params: &mcp.CallToolParams{
					Name:      tt.tool.name,
					Arguments: tt.args,
				},
			}

			// Execute handler
			result, err := handler(context.Background(), request)

			if (err != nil) != tt.wantError {
				t.Errorf("handler() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if result == nil {
				t.Error("handler() returned nil result")
				return
			}

			if result.IsError != tt.toolError {
				t.Errorf("result.IsError = %v, want %v", result.IsError, tt.toolError)
			}

			if len(result.Content) == 0 {
				t.Error("handler() returned empty content")
			}
		})
	}
}

func TestRegisterTools(t *testing.T) {
	// Create a test server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}, nil)

	// Create a test registry with mock tools
	registry := tools.NewRegistry()

	// Register test tools
	mockTool1 := &mockTool{
		name:        "tool1",
		description: "Test tool 1",
		schemaType:  "jsonschema",
	}
	mockTool2 := &mockTool{
		name:        "tool2",
		description: "Test tool 2",
		schemaType:  "legacy",
	}

	registry.Register(mockTool1)
	registry.Register(mockTool2)

	// Register tools with server
	registerTools(server, registry)

	// Verify tools were registered (check server has tools)
	// Note: The MCP server doesn't expose a way to check registered tools directly,
	// so we just verify the function completes without panic

	// Test with tool that has invalid schema
	registryWithBadTool := tools.NewRegistry()
	badTool := &badSchemaTool{
		mockTool: mockTool{
			name:        "bad-tool",
			description: "Bad tool",
		},
	}
	registryWithBadTool.Register(badTool)

	// This should not panic, just log the error
	registerTools(server, registryWithBadTool)
}

// badSchemaTool returns a schema that can't be marshalled
type badSchemaTool struct {
	mockTool
}

func (b *badSchemaTool) GetInputSchema() interface{} {
	// Return something that will fail JSON marshalling
	return make(chan int)
}
