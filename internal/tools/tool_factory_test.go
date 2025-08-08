package tools

import (
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/config"
	"github.com/raja-aiml/webex-mcp-server/internal/testutil"
)

func TestNewToolBase(t *testing.T) {
	schema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"test": StringProperty("Test field"),
		},
	}

	tool := NewToolBase("test-tool", "Test description", schema)

	if tool.Name() != "test-tool" {
		t.Errorf("Expected name 'test-tool', got %s", tool.Name())
	}

	if tool.Description() != "Test description" {
		t.Errorf("Expected description 'Test description', got %s", tool.Description())
	}

	if tool.GetInputSchema() != schema {
		t.Error("Expected schema to match input schema")
	}
}

func TestNewToolBaseWithConfig(t *testing.T) {
	cfg := &config.Config{
		WebexAPIKey:     "test-token",
		WebexAPIBaseURL: "https://api.test.com",
	}

	schema := &jsonschema.Schema{Type: "object"}

	tool := NewToolBaseWithConfig("test-tool", "Test description", schema, cfg)

	if tool.Name() != "test-tool" {
		t.Errorf("Expected name 'test-tool', got %s", tool.Name())
	}

	// Verify client is not yet initialized (lazy initialization)
	if tool.client != nil {
		t.Error("Expected client to be nil until first use (lazy initialization)")
	}

	// Test that ensureClient works
	if err := tool.ensureClient(); err != nil {
		t.Errorf("ensureClient() failed: %v", err)
	}

	// Now client should be initialized
	if tool.client == nil {
		t.Error("Expected client to be initialized after ensureClient()")
	}

	// Verify config
	if tool.config != cfg {
		t.Error("Expected config to be set")
	}
}

func TestToolBase_GetInputSchema(t *testing.T) {
	schema := &jsonschema.Schema{
		Type:        "object",
		Description: "Test schema",
	}

	tool := &ToolBase{
		name:        "test-tool",
		description: "Test description",
		schema:      schema,
	}

	result := tool.GetInputSchema()
	if result != schema {
		t.Error("GetInputSchema() should return the tool's schema")
	}
}

func TestSimpleSchema(t *testing.T) {
	description := "Test schema"
	properties := map[string]*jsonschema.Schema{
		"field1": StringProperty("Field 1"),
		"field2": IntegerProperty("Field 2"),
	}
	required := []string{"field1"}

	schema := SimpleSchema(description, properties, required)

	if schema.Type != "object" {
		t.Errorf("Expected type 'object', got %s", schema.Type)
	}

	if schema.Description != description {
		t.Errorf("Expected description '%s', got %s", description, schema.Description)
	}

	if len(schema.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(schema.Properties))
	}

	if len(schema.Required) != 1 {
		t.Errorf("Expected 1 required field, got %d", len(schema.Required))
	}

	if schema.Required[0] != "field1" {
		t.Errorf("Expected required field 'field1', got %s", schema.Required[0])
	}
}

func TestStringProperty(t *testing.T) {
	description := "Test string property"
	prop := StringProperty(description)

	if prop.Type != "string" {
		t.Errorf("Expected type 'string', got %s", prop.Type)
	}

	if prop.Description != description {
		t.Errorf("Expected description '%s', got %s", description, prop.Description)
	}
}

func TestIntegerProperty(t *testing.T) {
	description := "Test integer property"
	prop := IntegerProperty(description)

	if prop.Type != "integer" {
		t.Errorf("Expected type 'integer', got %s", prop.Type)
	}

	if prop.Description != description {
		t.Errorf("Expected description '%s', got %s", description, prop.Description)
	}
}

func TestBooleanProperty(t *testing.T) {
	description := "Test boolean property"
	prop := BooleanProperty(description)

	if prop.Type != "boolean" {
		t.Errorf("Expected type 'boolean', got %s", prop.Type)
	}

	if prop.Description != description {
		t.Errorf("Expected description '%s', got %s", description, prop.Description)
	}
}

func TestArrayProperty(t *testing.T) {
	description := "Test array property"
	items := StringProperty("Array item")
	prop := ArrayProperty(description, items)

	if prop.Type != "array" {
		t.Errorf("Expected type 'array', got %s", prop.Type)
	}

	if prop.Description != description {
		t.Errorf("Expected description '%s', got %s", description, prop.Description)
	}

	if prop.Items != items {
		t.Error("Expected items to match input")
	}
}

func TestObjectProperty(t *testing.T) {
	description := "Test object property"
	properties := map[string]*jsonschema.Schema{
		"nested": StringProperty("Nested field"),
	}
	prop := ObjectProperty(description, properties)

	if prop.Type != "object" {
		t.Errorf("Expected type 'object', got %s", prop.Type)
	}

	if prop.Description != description {
		t.Errorf("Expected description '%s', got %s", description, prop.Description)
	}

	if len(prop.Properties) != 1 {
		t.Errorf("Expected 1 property, got %d", len(prop.Properties))
	}
}

func TestPropertyHelpers_EdgeCases(t *testing.T) {
	t.Run("empty descriptions", func(t *testing.T) {
		prop := StringProperty("")
		if prop.Type != "string" {
			t.Error("Should still create valid property with empty description")
		}
	})

	t.Run("nested array", func(t *testing.T) {
		innerArray := ArrayProperty("Inner array", StringProperty("Item"))
		outerArray := ArrayProperty("Outer array", innerArray)

		if outerArray.Type != "array" {
			t.Error("Should create nested array properly")
		}
		if outerArray.Items.Type != "array" {
			t.Error("Inner array should be preserved")
		}
	})

	t.Run("object with properties", func(t *testing.T) {
		props := map[string]*jsonschema.Schema{
			"field1": StringProperty("Field 1"),
			"field2": IntegerProperty("Field 2"),
			"field3": BooleanProperty("Field 3"),
		}
		obj := ObjectProperty("Complex object", props)

		if len(obj.Properties) != 3 {
			t.Errorf("Expected 3 properties, got %d", len(obj.Properties))
		}
	})
}

func TestSimpleSchema_ComplexExample(t *testing.T) {
	// Test a complex schema like one used in real tools
	schema := SimpleSchema(
		"List rooms",
		map[string]*jsonschema.Schema{
			"teamId": StringProperty("List rooms associated with a team"),
			"type":   StringProperty("Room type: direct or group"),
			"sortBy": StringProperty("Sort results"),
			"max":    IntegerProperty("Limit the maximum number of rooms"),
		},
		[]string{}, // No required fields
	)

	if schema.Type != "object" {
		t.Error("Schema should be object type")
	}

	if len(schema.Properties) != 4 {
		t.Errorf("Expected 4 properties, got %d", len(schema.Properties))
	}

	if schema.Properties["max"].Type != "integer" {
		t.Error("'max' property should be integer type")
	}
}

func TestToolBase_WithNilClient(t *testing.T) {
	// Set up test environment with API key
	config.ResetForTesting()
	defaultClient = nil
	clientOnce = sync.Once{}
	clientErr = nil
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-token")
	defer func() {
		cleanup()
		config.ResetForTesting()
		defaultClient = nil
		clientOnce = sync.Once{}
		clientErr = nil
	}()

	tool := &ToolBase{
		name:        "test-tool",
		description: "Test",
		schema:      &jsonschema.Schema{Type: "object"},
		client:      nil,
	}

	// Ensure client initializes properly
	err := tool.ensureClient()
	if err != nil {
		t.Errorf("ensureClient() should handle nil client gracefully: %v", err)
	}

	if tool.client == nil {
		t.Error("Client should be initialized after ensureClient()")
	}
}

func TestSchemaBuilder_ChainableAPI(t *testing.T) {
	// Test that the helper functions can be chained nicely
	schema := SimpleSchema(
		"Complex schema",
		map[string]*jsonschema.Schema{
			"stringField": StringProperty("A string"),
			"numberField": IntegerProperty("A number"),
			"boolField":   BooleanProperty("A boolean"),
			"arrayField": ArrayProperty("An array",
				ObjectProperty("Array items", map[string]*jsonschema.Schema{
					"nested": StringProperty("Nested field"),
				}),
			),
		},
		[]string{"stringField"},
	)

	if schema.Type != "object" {
		t.Error("Should create valid schema through chaining")
	}

	if len(schema.Required) != 1 || schema.Required[0] != "stringField" {
		t.Error("Should preserve required fields")
	}

	arrayProp := schema.Properties["arrayField"]
	if arrayProp.Type != "array" {
		t.Error("Array property should be preserved")
	}

	if arrayProp.Items.Type != "object" {
		t.Error("Nested object in array should be preserved")
	}
}
