package tools

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/raja-aiml/webex-mcp-server-go/internal/config"
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
	mockProvider := &config.MockProvider{
		Token:   "test-token",
		BaseURL: "https://api.test.com",
	}
	
	schema := &jsonschema.Schema{Type: "object"}
	
	tool := NewToolBaseWithConfig("test-tool", "Test description", schema, mockProvider)
	
	if tool.Name() != "test-tool" {
		t.Errorf("Expected name 'test-tool', got %s", tool.Name())
	}
	
	// Verify client was created
	if tool.client == nil {
		t.Error("Expected client to be initialized")
	}
	
	// Verify config provider
	if tool.configProvider != mockProvider {
		t.Error("Expected config provider to be set")
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
		t.Error("Expected GetInputSchema to return the configured schema")
	}
}

func TestSimpleSchema(t *testing.T) {
	props := map[string]*jsonschema.Schema{
		"field1": StringProperty("Field 1"),
		"field2": IntegerProperty("Field 2"),
	}
	
	schema := SimpleSchema(props, []string{"field1"})
	
	// Check type
	if schema.Type != "object" {
		t.Errorf("Expected type 'object', got %v", schema.Type)
	}
	
	// Check properties
	if len(schema.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(schema.Properties))
	}
	
	// Check required
	if len(schema.Required) != 1 || schema.Required[0] != "field1" {
		t.Errorf("Expected required=['field1'], got %v", schema.Required)
	}
}

func TestStringProperty(t *testing.T) {
	prop := StringProperty("Test string")
	
	if prop.Type != "string" {
		t.Errorf("Expected type 'string', got %v", prop.Type)
	}
	
	if prop.Description != "Test string" {
		t.Errorf("Expected description 'Test string', got %v", prop.Description)
	}
}

func TestIntegerProperty(t *testing.T) {
	prop := IntegerProperty("Test integer")
	
	if prop.Type != "integer" {
		t.Errorf("Expected type 'integer', got %v", prop.Type)
	}
	
	if prop.Description != "Test integer" {
		t.Errorf("Expected description 'Test integer', got %v", prop.Description)
	}
}

func TestBooleanProperty(t *testing.T) {
	prop := BooleanProperty("Test boolean")
	
	if prop.Type != "boolean" {
		t.Errorf("Expected type 'boolean', got %v", prop.Type)
	}
	
	if prop.Description != "Test boolean" {
		t.Errorf("Expected description 'Test boolean', got %v", prop.Description)
	}
}

func TestArrayProperty(t *testing.T) {
	itemProp := StringProperty("Item")
	prop := ArrayProperty("Test array", itemProp)
	
	if prop.Type != "array" {
		t.Errorf("Expected type 'array', got %v", prop.Type)
	}
	
	if prop.Description != "Test array" {
		t.Errorf("Expected description 'Test array', got %v", prop.Description)
	}
	
	if prop.Items != itemProp {
		t.Error("Expected items to match the provided schema")
	}
}

func TestObjectProperty(t *testing.T) {
	prop := ObjectProperty("Test object")
	
	if prop.Type != "object" {
		t.Errorf("Expected type 'object', got %v", prop.Type)
	}
	
	if prop.Description != "Test object" {
		t.Errorf("Expected description 'Test object', got %v", prop.Description)
	}
}

func TestPropertyHelpers_EdgeCases(t *testing.T) {
	// Test with empty descriptions
	t.Run("empty descriptions", func(t *testing.T) {
		prop := StringProperty("")
		if prop.Description != "" {
			t.Errorf("Expected empty description, got %v", prop.Description)
		}
	})
	
	// Test nested array property
	t.Run("nested array", func(t *testing.T) {
		innerArray := ArrayProperty("Inner array", StringProperty("String item"))
		outerArray := ArrayProperty("Outer array", innerArray)
		
		if outerArray.Type != "array" {
			t.Errorf("Expected type 'array', got %v", outerArray.Type)
		}
		
		if outerArray.Items.Type != "array" {
			t.Errorf("Expected nested items type 'array', got %v", outerArray.Items.Type)
		}
	})
	
	// Test object with properties
	t.Run("object with properties", func(t *testing.T) {
		prop := ObjectProperty("User object")
		
		// Can add properties after creation
		prop.Properties = map[string]*jsonschema.Schema{
			"name": StringProperty("User name"),
			"age":  IntegerProperty("User age"),
		}
		
		if len(prop.Properties) != 2 {
			t.Errorf("Expected 2 properties, got %d", len(prop.Properties))
		}
	})
}

func TestSimpleSchema_ComplexExample(t *testing.T) {
	// Test a more complex schema with nested properties
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"name":   StringProperty("User name"),
		"age":    IntegerProperty("User age"),
		"active": BooleanProperty("Is active"),
		"tags":   ArrayProperty("User tags", StringProperty("Tag")),
		"address": &jsonschema.Schema{
			Type:        "object",
			Description: "User address",
			Properties: map[string]*jsonschema.Schema{
				"street": StringProperty("Street name"),
				"city":   StringProperty("City name"),
				"zip":    StringProperty("Zip code"),
			},
			Required: []string{"street", "city"},
		},
	}, []string{"name", "active"})
	
	// Verify schema structure
	if schema.Type != "object" {
		t.Errorf("Expected type 'object', got %v", schema.Type)
	}
	
	if len(schema.Properties) != 5 {
		t.Errorf("Expected 5 properties, got %d", len(schema.Properties))
	}
	
	// Check required fields
	if len(schema.Required) != 2 {
		t.Errorf("Expected 2 required fields, got %d", len(schema.Required))
	}
	
	// Verify address structure
	addressProp := schema.Properties["address"]
	if addressProp.Type != "object" {
		t.Errorf("Expected address type 'object', got %v", addressProp.Type)
	}
	
	if len(addressProp.Properties) != 3 {
		t.Errorf("Expected 3 address properties, got %d", len(addressProp.Properties))
	}
	
	if len(addressProp.Required) != 2 {
		t.Errorf("Expected 2 required address fields, got %d", len(addressProp.Required))
	}
}

func TestToolBase_WithNilClient(t *testing.T) {
	tool := &ToolBase{
		name:        "test-tool",
		description: "Test tool",
		schema:      &jsonschema.Schema{Type: "object"},
		client:      nil, // Explicitly nil client
	}
	
	// Should not panic
	if tool.Name() != "test-tool" {
		t.Errorf("Expected name 'test-tool', got %s", tool.Name())
	}
	
	if tool.Description() != "Test tool" {
		t.Errorf("Expected description 'Test tool', got %s", tool.Description())
	}
}

func TestSchemaBuilder_ChainableAPI(t *testing.T) {
	// Test that we can build schemas in a fluent way
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"users": ArrayProperty("List of users", &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"id":       IntegerProperty("User ID"),
				"name":     StringProperty("User name"),
				"email":    StringProperty("User email"),
				"verified": BooleanProperty("Email verified"),
				"roles":    ArrayProperty("User roles", StringProperty("Role name")),
			},
			Required: []string{"id", "name", "email"},
		}),
		"total": IntegerProperty("Total user count"),
		"page":  IntegerProperty("Current page"),
	}, []string{"users", "total"})
	
	// Verify the structure
	if schema.Type != "object" {
		t.Error("Expected root type to be object")
	}
	
	// Check users array property
	usersProp := schema.Properties["users"]
	if usersProp.Type != "array" {
		t.Error("Expected users to be array type")
	}
	
	// Check user object schema
	userSchema := usersProp.Items
	if userSchema.Type != "object" {
		t.Error("Expected user items to be object type")
	}
	
	if len(userSchema.Properties) != 5 {
		t.Errorf("Expected 5 user properties, got %d", len(userSchema.Properties))
	}
	
	// Check nested array
	rolesProp := userSchema.Properties["roles"]
	if rolesProp.Type != "array" {
		t.Error("Expected roles to be array type")
	}
	
	if rolesProp.Items.Type != "string" {
		t.Error("Expected role items to be string type")
	}
}