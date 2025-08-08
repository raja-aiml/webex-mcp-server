package tools

import (
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/config"
	"github.com/raja-aiml/webex-mcp-server/internal/webex"
)

// ToolBase provides common functionality for all tools
// implementing DRY principle by eliminating repetitive code
type ToolBase struct {
	name        string
	description string
	schema      *jsonschema.Schema
	client      webex.HTTPClient
	config      *config.Config
}

// NewToolBase creates a base tool with common functionality
func NewToolBase(name, description string, schema *jsonschema.Schema) ToolBase {
	cfg, _ := config.Load()
	return NewToolBaseWithConfig(name, description, schema, cfg)
}

// NewToolBaseWithConfig creates a base tool with dependency injection
func NewToolBaseWithConfig(name, description string, schema *jsonschema.Schema, cfg *config.Config) ToolBase {
	// Initialize the client lazily - it will be set when first used
	return ToolBase{
		name:        name,
		description: description,
		schema:      schema,
		config:      cfg,
	}
}

// ensureClient ensures the HTTP client is initialized
func (t *ToolBase) ensureClient() error {
	if t.client == nil {
		var err error
		if t.config != nil {
			t.client, err = webex.NewClientWithConfig(t.config)
		} else {
			t.client, err = getDefaultClient()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Name returns the tool name
func (t *ToolBase) Name() string { return t.name }

// Description returns the tool description
func (t *ToolBase) Description() string { return t.description }

// GetInputSchema returns the tool's input schema
func (t *ToolBase) GetInputSchema() interface{} { return t.schema }

// Factory functions for creating schemas in a consistent way

// SimpleSchema creates a simple schema with properties and required fields
func SimpleSchema(description string, properties map[string]*jsonschema.Schema, required []string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "object",
		Description: description,
		Properties:  properties,
		Required:    required,
	}
}

// StringProperty creates a string property schema
func StringProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Description: description,
	}
}

// IntegerProperty creates an integer property schema
func IntegerProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "integer",
		Description: description,
	}
}

// BooleanProperty creates a boolean property schema
func BooleanProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "boolean",
		Description: description,
	}
}

// ArrayProperty creates an array property schema
func ArrayProperty(description string, items *jsonschema.Schema) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "array",
		Description: description,
		Items:       items,
	}
}

// ObjectProperty creates an object property schema
func ObjectProperty(description string, properties map[string]*jsonschema.Schema) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "object",
		Description: description,
		Properties:  properties,
	}
}
