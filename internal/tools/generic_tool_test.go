package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/webex"
)

// Mock WebexClient for testing
type mockWebexClient struct {
	GetFunc    func(endpoint string, params map[string]string) (map[string]interface{}, error)
	PostFunc   func(endpoint string, data interface{}) (map[string]interface{}, error)
	PutFunc    func(endpoint string, data interface{}) (map[string]interface{}, error)
	DeleteFunc func(endpoint string) error
}

func (m *mockWebexClient) Get(endpoint string, params map[string]string) (map[string]interface{}, error) {
	if m.GetFunc != nil {
		return m.GetFunc(endpoint, params)
	}
	return nil, errors.New("not implemented")
}

func (m *mockWebexClient) Post(endpoint string, data interface{}) (map[string]interface{}, error) {
	if m.PostFunc != nil {
		return m.PostFunc(endpoint, data)
	}
	return nil, errors.New("not implemented")
}

func (m *mockWebexClient) Put(endpoint string, data interface{}) (map[string]interface{}, error) {
	if m.PutFunc != nil {
		return m.PutFunc(endpoint, data)
	}
	return nil, errors.New("not implemented")
}

func (m *mockWebexClient) Delete(endpoint string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(endpoint)
	}
	return errors.New("not implemented")
}

// Test structures for generic tool testing
type testListParams struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Filter string `json:"filter,omitempty"`
}

type testCreateParams struct {
	Name        string `json:"name" jsonschema:"required"`
	Description string `json:"description,omitempty"`
}

type testIDParams struct {
	ID string `json:"id" jsonschema:"required"`
}

func TestNewGenericTool(t *testing.T) {
	schema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"name": {Type: "string", Description: "Item name"},
		},
		Required: []string{"name"},
	}

	executor := func(params *testCreateParams, client webex.HTTPClient) (interface{}, error) {
		return map[string]interface{}{"created": true}, nil
	}

	tool := NewGenericTool("test-tool", "Test tool", schema, executor)

	if tool.Name() != "test-tool" {
		t.Errorf("Expected name 'test-tool', got %s", tool.Name())
	}

	if tool.Description() != "Test tool" {
		t.Errorf("Expected description 'Test tool', got %s", tool.Description())
	}
}

func TestGenericTool_Execute(t *testing.T) {
	tests := []struct {
		name       string
		params     interface{}
		executor   interface{}
		mockFunc   func(*mockWebexClient)
		wantResult interface{}
		wantErr    bool
	}{
		{
			name:   "successful list execution",
			params: testListParams{Limit: 10, Filter: "active"},
			executor: func(params *testListParams, client webex.HTTPClient) (interface{}, error) {
				queryParams := map[string]string{
					"limit":  fmt.Sprintf("%d", params.Limit),
					"filter": params.Filter,
				}
				return client.Get("/items", queryParams)
			},
			mockFunc: func(m *mockWebexClient) {
				m.GetFunc = func(endpoint string, params map[string]string) (map[string]interface{}, error) {
					if endpoint != "/items" {
						t.Errorf("Expected endpoint /items, got %s", endpoint)
					}
					if params["limit"] != "10" {
						t.Errorf("Expected limit=10, got %s", params["limit"])
					}
					return map[string]interface{}{"items": []interface{}{}}, nil
				}
			},
			wantResult: map[string]interface{}{"items": []interface{}{}},
			wantErr:    false,
		},
		{
			name:   "successful create execution",
			params: testCreateParams{Name: "Test Item", Description: "A test item"},
			executor: func(params *testCreateParams, client webex.HTTPClient) (interface{}, error) {
				body := map[string]interface{}{
					"name":        params.Name,
					"description": params.Description,
				}
				return client.Post("/items", body)
			},
			mockFunc: func(m *mockWebexClient) {
				m.PostFunc = func(endpoint string, data interface{}) (map[string]interface{}, error) {
					body, ok := data.(map[string]interface{})
					if !ok {
						t.Errorf("Expected data to be map[string]interface{}, got %T", data)
						return nil, errors.New("invalid data type")
					}
					if body["name"] != "Test Item" {
						t.Errorf("Expected name='Test Item', got %v", body["name"])
					}
					return map[string]interface{}{"id": "123", "name": "Test Item"}, nil
				}
			},
			wantResult: map[string]interface{}{"id": "123", "name": "Test Item"},
			wantErr:    false,
		},
		// Validation now relies on JSON schema, not reflection
		// This test case is removed as validation happens at MCP level
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockWebexClient{}
			tt.mockFunc(mockClient)

			// Create schema based on params type
			schema := &jsonschema.Schema{Type: "object"}

			// Create tool based on params type
			var result interface{}
			var err error

			switch p := tt.params.(type) {
			case testListParams:
				tool := &GenericTool[testListParams]{
					ToolBase: ToolBase{
						name:        "test-tool",
						description: "Test tool",
						schema:      schema,
						client:      mockClient,
					},
					executor: tt.executor.(func(*testListParams, webex.HTTPClient) (interface{}, error)),
				}
				jsonArgs, _ := json.Marshal(p)
				result, err = tool.Execute(jsonArgs)

			case testCreateParams:
				tool := &GenericTool[testCreateParams]{
					ToolBase: ToolBase{
						name:        "test-tool",
						description: "Test tool",
						schema:      schema,
						client:      mockClient,
					},
					executor: tt.executor.(func(*testCreateParams, webex.HTTPClient) (interface{}, error)),
				}
				jsonArgs, _ := json.Marshal(p)
				result, err = tool.Execute(jsonArgs)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != nil {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Expected result to be map[string]interface{}, got %T", result)
					return
				}
				wantMap := tt.wantResult.(map[string]interface{})
				for k, v := range wantMap {
					// Compare values - handle slice comparison
					switch expected := v.(type) {
					case []interface{}:
						if actual, ok := resultMap[k].([]interface{}); ok {
							if len(actual) != len(expected) {
								t.Errorf("Result[%s] length = %d, want %d", k, len(actual), len(expected))
							}
						} else {
							t.Errorf("Result[%s] is not a slice", k)
						}
					default:
						if resultMap[k] != v {
							t.Errorf("Result[%s] = %v, want %v", k, resultMap[k], v)
						}
					}
				}
			}
		})
	}
}

func TestGenericTool_ExecuteWithMap(t *testing.T) {
	mockClient := &mockWebexClient{
		GetFunc: func(endpoint string, params map[string]string) (map[string]interface{}, error) {
			return map[string]interface{}{"result": "success"}, nil
		},
	}

	schema := &jsonschema.Schema{Type: "object"}
	tool := &GenericTool[testListParams]{
		ToolBase: ToolBase{
			name:        "test-tool",
			description: "Test tool",
			schema:      schema,
			client:      mockClient,
		},
		executor: func(params *testListParams, client webex.HTTPClient) (interface{}, error) {
			return client.Get("/test", nil)
		},
	}

	result, err := tool.ExecuteWithMap(map[string]interface{}{"limit": 10})
	if err != nil {
		t.Fatalf("ExecuteWithMap() error = %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result to be map[string]interface{}, got %T", result)
	}

	if resultMap["result"] != "success" {
		t.Errorf("Expected result=success, got %v", resultMap["result"])
	}
}

// TestValidateRequired removed - validation now relies on JSON schema
/*
func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name    string
		params  interface{}
		wantErr bool
	}{
		{
			name: "all required fields present",
			params: &struct {
				Required1 string `json:"required1" required:"true"`
				Required2 int    `json:"required2" required:"true"`
				Optional  string `json:"optional"`
			}{
				Required1: "value1",
				Required2: 42,
				Optional:  "",
			},
			wantErr: false,
		},
		{
			name: "missing required string field",
			params: &struct {
				Required string `json:"required" required:"true"`
			}{
				Required: "",
			},
			wantErr: true,
		},
		{
			name: "zero int value is invalid",
			params: &struct {
				Required int `json:"required" required:"true"`
			}{
				Required: 0,
			},
			wantErr: true,  // isZeroValue treats 0 as zero value
		},
		{
			name: "false bool value is invalid",
			params: &struct {
				Required bool `json:"required" required:"true"`
			}{
				Required: false,
			},
			wantErr: true,  // isZeroValue treats false as zero value
		},
		{
			name: "nil pointer field is invalid",
			params: &struct {
				Required *string `json:"required" required:"true"`
			}{
				Required: nil,
			},
			wantErr: true,
		},
		{
			name: "empty slice is invalid",
			params: &struct {
				Required []string `json:"required" required:"true"`
			}{
				Required: []string{},
			},
			wantErr: true,
		},
		{
			name: "non-empty slice is valid",
			params: &struct {
				Required []string `json:"required" required:"true"`
			}{
				Required: []string{"item"},
			},
			wantErr: false,
		},
		{
			name: "fields without required tag are not validated",
			params: &struct {
				NotRequired string `json:"notRequired"`
			}{
				NotRequired: "",
			},
			wantErr: false,
		},
		{
			name: "jsonschema required tag is not used",
			params: &struct {
				Required string `json:"required" jsonschema:"required"`
			}{
				Required: "",
			},
			wantErr: false,  // Only required:"true" tag is checked
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequired(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
*/

// TestIsZeroValue removed - no longer using reflection-based validation
/*
func TestIsZeroValue(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{"empty string", "", true},
		{"non-empty string", "test", false},
		{"zero int", 0, true},        // isZeroValue considers 0 as zero
		{"non-zero int", 42, false},
		{"false bool", false, true},   // isZeroValue considers false as zero (!v.Bool())
		{"true bool", true, false},
		{"nil interface", nil, true},
		{"nil pointer", (*string)(nil), true},
		{"empty slice", []string{}, true},
		{"non-empty slice", []string{"a"}, false},
		{"empty map", map[string]interface{}{}, true},
		{"non-empty map", map[string]interface{}{"a": 1}, false},
		{"zero float", 0.0, true},     // isZeroValue considers 0.0 as zero
		{"non-zero float", 3.14, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handle nil separately
			if tt.value == nil {
				// isZeroValue expects a reflect.Value, not a nil interface
				// For nil values, we should consider them as zero
				if !tt.want {
					t.Errorf("Expected nil to be zero value")
				}
				return
			}

			v := reflect.ValueOf(tt.value)
			if got := isZeroValue(v); got != tt.want {
				t.Errorf("isZeroValue(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
*/

func TestQueryParams(t *testing.T) {
	// Test with map[string]interface{}
	params := map[string]interface{}{
		"param1": "value1",
		"param2": "", // Empty string should be ignored
		"param3": 123,
		"param4": true,
		"param5": nil,   // nil should be ignored
		"param6": 0,     // Zero should be ignored
		"param7": false, // False should be included
	}

	queryParams := QueryParams(params)

	if queryParams["param1"] != "value1" {
		t.Errorf("Expected param1=value1, got %s", queryParams["param1"])
	}

	if _, exists := queryParams["param2"]; exists {
		t.Error("Expected param2 to be excluded (empty string)")
	}

	if queryParams["param3"] != "123" {
		t.Errorf("Expected param3=123, got %s", queryParams["param3"])
	}

	if queryParams["param4"] != "true" {
		t.Errorf("Expected param4=true, got %s", queryParams["param4"])
	}

	if _, exists := queryParams["param5"]; exists {
		t.Error("Expected param5 to be excluded (nil value)")
	}

	if _, exists := queryParams["param6"]; exists {
		t.Error("Expected param6 to be excluded (zero value)")
	}

	if queryParams["param7"] != "false" {
		t.Errorf("Expected param7=false, got %s", queryParams["param7"])
	}
}

func TestNewListTool(t *testing.T) {
	properties := map[string]*jsonschema.Schema{
		"limit":  IntegerProperty("Maximum number of items"),
		"offset": IntegerProperty("Number of items to skip"),
	}

	tool := NewListTool[testListParams]("test-list", "List test items", "/items", properties, []string{})

	if tool.Name() != "test-list" {
		t.Errorf("Expected name 'test-list', got %s", tool.Name())
	}

	if tool.Description() != "List test items" {
		t.Errorf("Expected description 'List test items', got %s", tool.Description())
	}
}

func TestNewGetTool(t *testing.T) {
	tool := NewGetTool("test-get", "Get test item", "/items", "id", "The ID of the item")

	if tool.Name() != "test-get" {
		t.Errorf("Expected name 'test-get', got %s", tool.Name())
	}

	if tool.Description() != "Get test item" {
		t.Errorf("Expected description 'Get test item', got %s", tool.Description())
	}
}

func TestNewDeleteTool(t *testing.T) {
	tool := NewDeleteTool("test-delete", "Delete test item", "/items", "id", "The ID of the item")

	if tool.Name() != "test-delete" {
		t.Errorf("Expected name 'test-delete', got %s", tool.Name())
	}

	if tool.Description() != "Delete test item" {
		t.Errorf("Expected description 'Delete test item', got %s", tool.Description())
	}
}

func TestGenericTool_GetInputSchema(t *testing.T) {
	schema := &jsonschema.Schema{
		Type:        "object",
		Description: "Test schema",
		Properties: map[string]*jsonschema.Schema{
			"field1": StringProperty("Field 1"),
		},
		Required: []string{"field1"},
	}

	tool := &GenericTool[testCreateParams]{
		ToolBase: ToolBase{
			name:        "test-tool",
			description: "Test tool",
			schema:      schema,
		},
	}

	resultSchema := tool.GetInputSchema()

	if resultSchema == nil {
		t.Fatal("Expected non-nil schema")
	}

	// Verify it returns the same schema
	if resultSchema != schema {
		t.Error("Expected GetInputSchema to return the configured schema")
	}
}

// Test that GenericTool implements the Tool interface
func TestGenericTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = &GenericTool[testCreateParams]{}

	// This test will fail to compile if GenericTool doesn't implement Tool interface
	t.Log("GenericTool implements Tool interface")
}
