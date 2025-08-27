package tools

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/config"
	"github.com/raja-aiml/webex-mcp-server/internal/webex"
)

// ToolBase provides common functionality for all tools
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

func (t *ToolBase) Name() string                { return t.name }
func (t *ToolBase) Description() string         { return t.description }
func (t *ToolBase) GetInputSchema() interface{} { return t.schema }

// SimpleSchema creates a simple schema with properties and required fields
func SimpleSchema(description string, properties map[string]*jsonschema.Schema, required []string) *jsonschema.Schema {
	schema := &jsonschema.Schema{
		Type:        "object",
		Description: description,
		Properties:  properties,
	}

	if len(required) > 0 {
		schema.Required = required
	}

	return schema
}

// Property factory functions
func StringProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Description: description,
	}
}

// RequiredStringProperty creates a required string property
func RequiredStringProperty(description string) *jsonschema.Schema {
	minLen := 1
	return &jsonschema.Schema{
		Type:        "string",
		Description: description,
		MinLength:   &minLen,
	}
}

func IntegerProperty(description string) *jsonschema.Schema {
	min := 0.0
	return &jsonschema.Schema{
		Type:        "integer",
		Description: description,
		Minimum:     &min,
	}
}

func BooleanProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "boolean",
		Description: description,
	}
}

func ArrayProperty(description string, items *jsonschema.Schema) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "array",
		Description: description,
		Items:       items,
	}
}

func ObjectProperty(description string, properties map[string]*jsonschema.Schema) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "object",
		Description: description,
		Properties:  properties,
	}
}
