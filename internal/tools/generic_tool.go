package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
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
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Validate required fields using reflection
	if err := validateRequired(&params); err != nil {
		return nil, err
	}

	if err := t.ensureClient(); err != nil {
		return nil, fmt.Errorf("failed to initialize client: %w", err)
	}
	return t.executor(&params, t.client)
}

// ExecuteWithMap implements the Tool interface
func (t *GenericTool[T]) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}

// validateRequired checks that required fields are not empty
func validateRequired(params interface{}) error {
	v := reflect.ValueOf(params).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Check if field has "required" tag
		if required := field.Tag.Get("required"); required == "true" {
			if isZeroValue(value) {
				jsonTag := field.Tag.Get("json")
				if jsonTag == "" {
					jsonTag = field.Name
				}
				return fmt.Errorf("%s is required", jsonTag)
			}
		}
	}

	return nil
}

// isZeroValue checks if a reflect.Value is a zero value
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Map:
		return v.IsNil() || v.Len() == 0
	default:
		return v.IsZero()
	}
}

// QueryParams helps build query parameters from a struct
func QueryParams(params interface{}) map[string]string {
	if params == nil {
		return nil
	}
	
	result := make(map[string]string)
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	
	// Only process struct types
	if v.Kind() != reflect.Struct {
		return nil
	}
	
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip zero values unless explicitly marked to include
		if isZeroValue(value) && field.Tag.Get("includeZero") != "true" {
			continue
		}

		// Get the query parameter name from the tag
		paramName := field.Tag.Get("query")
		if paramName == "" {
			paramName = field.Tag.Get("json")
			if paramName == "" {
				paramName = field.Name
			}
		}

		// Convert value to string
		switch value.Kind() {
		case reflect.String:
			if s := value.String(); s != "" {
				result[paramName] = s
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Int() != 0 {
				result[paramName] = strconv.FormatInt(value.Int(), 10)
			}
		case reflect.Bool:
			if value.Bool() {
				result[paramName] = "true"
			}
		}
	}

	return result
}

// --- Helper Types for Common Patterns ---

// ListParams represents common parameters for list operations
type ListParams struct {
	Max int `json:"max,omitempty" query:"max"`
}

// IDParams represents operations that require an ID
type IDParams struct {
	ID string `json:"id" required:"true"`
}

// --- Factory Functions for Common Tool Patterns ---

// NewListTool creates a generic list tool
func NewListTool[T any](name, description, endpoint string, properties map[string]*jsonschema.Schema) *GenericTool[T] {
	schema := SimpleSchema("List items from the API endpoint.", properties, []string{})

	return NewGenericTool(name, description, schema, func(params *T, client webex.HTTPClient) (interface{}, error) {
		queryParams := QueryParams(params)
		return client.Get(endpoint, queryParams)
	})
}

// NewGetTool creates a generic get-by-id tool
func NewGetTool(name, description, endpoint, idField, idDescription string) Tool {
	schema := SimpleSchema("Get a specific item by ID.", map[string]*jsonschema.Schema{
		idField: StringProperty(idDescription),
	}, []string{idField})

	return NewGenericTool(name, description, schema, func(params *map[string]interface{}, client webex.HTTPClient) (interface{}, error) {
		id, ok := (*params)[idField].(string)
		if !ok || id == "" {
			return nil, fmt.Errorf("%s is required", idField)
		}
		return client.Get(fmt.Sprintf("%s/%s", endpoint, id), nil)
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
	// Add ID field to properties
	allProperties := make(map[string]*jsonschema.Schema)
	for k, v := range properties {
		allProperties[k] = v
	}
	allProperties[idField] = StringProperty("The ID of the item to update")

	// Add ID field to required
	allRequired := append([]string{idField}, required...)

	schema := SimpleSchema("Update an existing item.", allProperties, allRequired)

	return NewGenericTool(name, description, schema, func(params *T, client webex.HTTPClient) (interface{}, error) {
		// Extract ID using reflection
		v := reflect.ValueOf(params).Elem()
		idValue := ""

		// Find the ID field
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if field.Tag.Get("json") == idField {
				idValue = v.Field(i).String()
				break
			}
		}

		if idValue == "" {
			return nil, fmt.Errorf("%s is required", idField)
		}

		return client.Put(fmt.Sprintf("%s/%s", endpoint, idValue), params)
	})
}

// NewDeleteTool creates a generic delete tool
func NewDeleteTool(name, description, endpoint, idField, idDescription string) Tool {
	schema := SimpleSchema("Delete an item by ID.", map[string]*jsonschema.Schema{
		idField: StringProperty(idDescription),
	}, []string{idField})

	return NewGenericTool(name, description, schema, func(params *map[string]interface{}, client webex.HTTPClient) (interface{}, error) {
		id, ok := (*params)[idField].(string)
		if !ok || id == "" {
			return nil, fmt.Errorf("%s is required", idField)
		}
		err := client.Delete(fmt.Sprintf("%s/%s", endpoint, id))
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{"success": true}, nil
	})
}
