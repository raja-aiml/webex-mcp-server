package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/webex"
)

// GenericTool provides a generic implementation for tools
type GenericTool[T any] struct {
	ToolBase
	executor func(*T, webex.HTTPClient) (interface{}, error)
}

// NewGenericTool creates a new generic tool
func NewGenericTool[T any](name, description string, schema *jsonschema.Schema, executor func(*T, webex.HTTPClient) (interface{}, error)) *GenericTool[T] {
	return &GenericTool[T]{
		ToolBase: NewToolBase(name, description, schema),
		executor: executor,
	}
}

// Execute implements the Tool interface
func (t *GenericTool[T]) Execute(args json.RawMessage) (interface{}, error) {
	var params T
	if err := json.Unmarshal(args, &params); err != nil {
		// Provide more helpful error message
		return nil, fmt.Errorf("invalid arguments format: %w. Please check the tool schema for required fields", err)
	}

	if err := t.ensureClient(); err != nil {
		return nil, fmt.Errorf("service initialization failed: %w. Please check your API credentials", err)
	}

	result, err := t.executor(&params, t.client)
	if err != nil {
		// Wrap errors with more context
		return nil, fmt.Errorf("%s failed: %w", t.name, err)
	}

	return result, nil
}

// ExecuteWithMap implements the Tool interface
func (t *GenericTool[T]) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// SimpleTool provides a simple implementation without generics for map-based operations
type SimpleTool struct {
	ToolBase
	executor func(map[string]interface{}, webex.HTTPClient) (interface{}, error)
}

// NewSimpleTool creates a new simple tool
func NewSimpleTool(name, description string, schema *jsonschema.Schema, executor func(map[string]interface{}, webex.HTTPClient) (interface{}, error)) *SimpleTool {
	return &SimpleTool{
		ToolBase: NewToolBase(name, description, schema),
		executor: executor,
	}
}

// Execute implements the Tool interface
func (t *SimpleTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if len(args) > 0 {
		if err := json.Unmarshal(args, &params); err != nil {
			return nil, fmt.Errorf("failed to parse arguments: %w", err)
		}
	}

	if err := t.ensureClient(); err != nil {
		return nil, fmt.Errorf("failed to initialize client: %w", err)
	}
	return t.executor(params, t.client)
}

// ExecuteWithMap implements the Tool interface
func (t *SimpleTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// --- Helper Functions ---

// QueryParams converts a map to query parameters
func QueryParams(params interface{}) map[string]string {
	if params == nil {
		return nil
	}

	// Handle map[string]interface{} directly
	if m, ok := params.(map[string]interface{}); ok {
		return mapToQueryParams(m)
	}

	// Handle pointer to map
	if m, ok := params.(*map[string]interface{}); ok && m != nil {
		return mapToQueryParams(*m)
	}

	// For other types, return empty map
	return make(map[string]string)
}

// mapToQueryParams converts a map to query parameters
func mapToQueryParams(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if v != nil {
			switch val := v.(type) {
			case string:
				if val != "" {
					result[k] = val
				}
			case int:
				if val != 0 {
					result[k] = fmt.Sprintf("%d", val)
				}
			case int8, int16, int32, int64:
				result[k] = fmt.Sprintf("%d", val)
			case uint, uint8, uint16, uint32, uint64:
				result[k] = fmt.Sprintf("%d", val)
			case float32, float64:
				result[k] = strconv.FormatFloat(val.(float64), 'f', -1, 64)
			case bool:
				result[k] = strconv.FormatBool(val)
			default:
				// For any other type, use fmt.Sprintf
				s := fmt.Sprintf("%v", val)
				if s != "" && s != "<nil>" {
					result[k] = s
				}
			}
		}
	}
	return result
}

// --- Helper Types for Common Patterns ---

// ListParams represents common parameters for list operations
type ListParams struct {
	Max int `json:"max,omitempty"`
}

// IDParams represents operations that require an ID
type IDParams struct {
	ID string `json:"id"`
}

// --- Factory Functions for Common Tool Patterns ---

// NewListTool creates a generic list tool
func NewListTool[T any](name, description, endpoint string, properties map[string]*jsonschema.Schema, required []string) *GenericTool[T] {
	schema := SimpleSchema("List items from the API endpoint.", properties, required)

	return NewGenericTool(name, description, schema, func(params *T, client webex.HTTPClient) (interface{}, error) {
		// Convert params to map for query parameters
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}

		var paramsMap map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &paramsMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal to map: %w", err)
		}

		queryParams := mapToQueryParams(paramsMap)
		return client.Get(endpoint, queryParams)
	})
}

// NewGetTool creates a generic get-by-id tool
func NewGetTool(name, description, endpoint, idField, idDescription string) Tool {
	schema := SimpleSchema("Get a specific item by ID.", map[string]*jsonschema.Schema{
		idField: StringProperty(idDescription),
	}, []string{idField})

	return NewSimpleTool(name, description, schema, func(params map[string]interface{}, client webex.HTTPClient) (interface{}, error) {
		id, ok := params[idField]
		if !ok || id == nil {
			return nil, fmt.Errorf("%s is required", idField)
		}

		idStr := fmt.Sprintf("%v", id)
		if idStr == "" || idStr == "<nil>" {
			return nil, fmt.Errorf("%s cannot be empty", idField)
		}

		return client.Get(fmt.Sprintf("%s/%s", endpoint, idStr), nil)
	})
}

// NewCreateTool creates a generic create tool
func NewCreateTool[T any](name, description, endpoint string, properties map[string]*jsonschema.Schema, required []string) *GenericTool[T] {
	schema := SimpleSchema("Create a new item.", properties, required)

	return NewGenericTool(name, description, schema, func(params *T, client webex.HTTPClient) (interface{}, error) {
		return client.Post(endpoint, params)
	})
}

// NewUpdateTool creates a generic update tool
func NewUpdateTool[T any](name, description, endpoint, idField string, properties map[string]*jsonschema.Schema, required []string) *GenericTool[T] {
	// Create a copy of properties to avoid mutation
	allProperties := make(map[string]*jsonschema.Schema, len(properties)+1)
	for k, v := range properties {
		allProperties[k] = v
	}
	allProperties[idField] = StringProperty("The ID of the item to update")

	// Add ID field to required if not already present
	allRequired := make([]string, 0, len(required)+1)
	hasIdField := false
	for _, field := range required {
		if field == idField {
			hasIdField = true
		}
		allRequired = append(allRequired, field)
	}
	if !hasIdField {
		allRequired = append([]string{idField}, allRequired...)
	}

	schema := SimpleSchema("Update an existing item.", allProperties, allRequired)

	return NewGenericTool(name, description, schema, func(params *T, client webex.HTTPClient) (interface{}, error) {
		// Convert params to map to extract ID
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %w", err)
		}

		var paramsMap map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &paramsMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal to map: %w", err)
		}

		id, ok := paramsMap[idField]
		if !ok || id == nil {
			return nil, fmt.Errorf("%s is required", idField)
		}

		idStr := fmt.Sprintf("%v", id)
		if idStr == "" || idStr == "<nil>" {
			return nil, fmt.Errorf("%s cannot be empty", idField)
		}

		return client.Put(fmt.Sprintf("%s/%s", endpoint, idStr), params)
	})
}

// NewDeleteTool creates a generic delete tool
func NewDeleteTool(name, description, endpoint, idField, idDescription string) Tool {
	schema := SimpleSchema("Delete an item by ID.", map[string]*jsonschema.Schema{
		idField: StringProperty(idDescription),
	}, []string{idField})

	return NewSimpleTool(name, description, schema, func(params map[string]interface{}, client webex.HTTPClient) (interface{}, error) {
		id, ok := params[idField]
		if !ok || id == nil {
			return nil, fmt.Errorf("%s is required", idField)
		}

		idStr := fmt.Sprintf("%v", id)
		if idStr == "" || idStr == "<nil>" {
			return nil, fmt.Errorf("%s cannot be empty", idField)
		}

		err := client.Delete(fmt.Sprintf("%s/%s", endpoint, idStr))
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{"success": true}, nil
	})
}
