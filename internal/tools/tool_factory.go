package tools

import (
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/raja-aiml/webex-mcp-server-go/internal/config"
	"github.com/raja-aiml/webex-mcp-server-go/internal/webex"
)

// ToolBase provides common functionality for all tools
// implementing DRY principle by eliminating repetitive code
type ToolBase struct {
	name           string
	description    string
	schema         *jsonschema.Schema
	client         webex.HTTPClient
	configProvider config.Provider
}

// NewToolBase creates a base tool with common functionality
func NewToolBase(name, description string, schema *jsonschema.Schema) ToolBase {
	return NewToolBaseWithConfig(name, description, schema, config.NewDefaultProvider())
}

// NewToolBaseWithConfig creates a base tool with dependency injection
func NewToolBaseWithConfig(name, description string, schema *jsonschema.Schema, configProvider config.Provider) ToolBase {
	return ToolBase{
		name:           name,
		description:    description,
		schema:         schema,
		client:         webex.NewClientWithConfig(configProvider),
		configProvider: configProvider,
	}
}

// Name returns the tool name
func (t *ToolBase) Name() string { return t.name }

// Description returns the tool description
func (t *ToolBase) Description() string { return t.description }

// GetInputSchema returns the JSON schema
func (t *ToolBase) GetInputSchema() interface{} { return t.schema }

// Note: ExecuteWithMap is not implemented here because each tool that embeds ToolBase
// should implement it by calling ExecuteWithMapBase(tool, args) from base.go

// SimpleSchema creates a basic object schema with properties
func SimpleSchema(properties map[string]*jsonschema.Schema, required []string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
}

// StringProperty creates a string schema property
func StringProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Description: description,
	}
}

// IntegerProperty creates an integer schema property
func IntegerProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "integer",
		Description: description,
	}
}

// BooleanProperty creates a boolean schema property
func BooleanProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "boolean",
		Description: description,
	}
}

// ArrayProperty creates an array schema property
func ArrayProperty(description string, items *jsonschema.Schema) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "array",
		Description: description,
		Items:       items,
	}
}

// ObjectProperty creates an object schema property
func ObjectProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "object",
		Description: description,
	}
}
