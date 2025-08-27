package tools

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/config"
	"github.com/raja-aiml/webex-mcp-server/internal/webex"
)

// ToolBase provides common functionality for all tools - TypeScript base class pattern
// implementing DRY principle by eliminating repetitive code
type ToolBase struct {
	name        string                 // Tool name - private field for interface compliance
	description string                 // Tool description - private field for interface compliance
	schema      *jsonschema.Schema     // Tool input schema - private field
	client      webex.HTTPClient       // HTTP client - private field, excluded from JSON
	config      *config.Config         // Configuration - private field, excluded from JSON
	metadata    map[string]interface{} // TypeScript metadata pattern - private field
}

// ToolConfig represents tool configuration - TypeScript config interface pattern
type ToolConfig struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      *jsonschema.Schema     `json:"schema"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ToolFactoryOptions represents tool factory options - TypeScript options pattern
type ToolFactoryOptions struct {
	Config   *config.Config         `json:"config,omitempty"`
	Client   webex.HTTPClient       `json:"-"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewToolBase creates a base tool with common functionality - TypeScript constructor pattern
func NewToolBase(name, description string, schema *jsonschema.Schema) ToolBase {
	cfg, _ := config.Load()
	return NewToolBaseWithConfig(name, description, schema, cfg)
}

// NewToolBaseWithConfig creates a base tool with dependency injection - TypeScript DI pattern
func NewToolBaseWithConfig(name, description string, schema *jsonschema.Schema, cfg *config.Config) ToolBase {
	return ToolBase{
		name:        name,
		description: description,
		schema:      schema,
		config:      cfg,
		metadata:    make(map[string]interface{}), // TypeScript empty object pattern
	}
}

// NewToolBaseWithOptions creates a base tool with options - TypeScript options pattern
func NewToolBaseWithOptions(toolConfig ToolConfig, options ToolFactoryOptions) ToolBase {
	tool := ToolBase{
		name:        toolConfig.Name,
		description: toolConfig.Description,
		schema:      toolConfig.Schema,
		config:      options.Config,
		client:      options.Client,
		metadata:    make(map[string]interface{}),
	}

	// Merge metadata - TypeScript object spread equivalent
	if toolConfig.Metadata != nil {
		for key, value := range toolConfig.Metadata {
			tool.metadata[key] = value
		}
	}
	if options.Metadata != nil {
		for key, value := range options.Metadata {
			tool.metadata[key] = value
		}
	}

	return tool
}

// ensureClient ensures the HTTP client is initialized - TypeScript lazy initialization pattern
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

// Name returns the tool name - TypeScript getter pattern (interface compliance)
func (t *ToolBase) Name() string { return t.name }

// Description returns the tool description - TypeScript getter pattern (interface compliance)
func (t *ToolBase) Description() string { return t.description }

// GetInputSchema returns the tool's input schema - TypeScript getter pattern
func (t *ToolBase) GetInputSchema() interface{} { return t.schema }

// GetMetadata returns tool metadata - TypeScript getter pattern
func (t *ToolBase) GetMetadata() map[string]interface{} {
	if t.metadata == nil {
		return make(map[string]interface{})
	}
	return t.metadata
}

// SetMetadata sets tool metadata - TypeScript setter pattern
func (t *ToolBase) SetMetadata(key string, value interface{}) {
	if t.metadata == nil {
		t.metadata = make(map[string]interface{})
	}
	t.metadata[key] = value
}

// ToJSON returns tool as JSON-serializable object - TypeScript toJSON pattern
func (t *ToolBase) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"name":        t.name,
		"description": t.description,
		"schema":      t.schema,
		"metadata":    t.metadata,
	}
}

// Schema factory functions for creating schemas in a TypeScript-compliant way

// SchemaBuilder provides fluent API for schema creation - TypeScript builder pattern
type SchemaBuilder struct {
	schema *jsonschema.Schema
}

// NewSchemaBuilder creates a new schema builder - TypeScript constructor pattern
func NewSchemaBuilder(schemaType, description string) *SchemaBuilder {
	return &SchemaBuilder{
		schema: &jsonschema.Schema{
			Type:        schemaType,
			Description: description,
			Properties:  make(map[string]*jsonschema.Schema),
		},
	}
}

// AddProperty adds a property to the schema - TypeScript fluent API pattern
func (sb *SchemaBuilder) AddProperty(name string, property *jsonschema.Schema) *SchemaBuilder {
	if sb.schema.Properties == nil {
		sb.schema.Properties = make(map[string]*jsonschema.Schema)
	}
	sb.schema.Properties[name] = property
	return sb
}

// SetRequired sets required fields - TypeScript fluent API pattern
func (sb *SchemaBuilder) SetRequired(required []string) *SchemaBuilder {
	sb.schema.Required = required
	return sb
}

// Build returns the built schema - TypeScript builder pattern
func (sb *SchemaBuilder) Build() *jsonschema.Schema {
	// Set additionalProperties to false for strict validation - TypeScript strict mode
	if sb.schema.AdditionalProperties == nil {
		schema := &jsonschema.Schema{Type: "boolean"}
		sb.schema.AdditionalProperties = schema
	}
	return sb.schema
}

// SimpleSchema creates a simple schema with properties and required fields - TypeScript factory pattern
func SimpleSchema(description string, properties map[string]*jsonschema.Schema, required []string) *jsonschema.Schema {
	schema := &jsonschema.Schema{
		Type:        "object",
		Description: description,
		Properties:  properties,
	}

	// Only set required if there are actually required fields
	if len(required) > 0 {
		schema.Required = required
	}

	// Set additionalProperties to false for strict validation
	additionalProps := &jsonschema.Schema{Type: "boolean"}
	schema.AdditionalProperties = additionalProps

	return schema
}

// Property factory functions with TypeScript naming conventions

// StringProperty creates a string property schema - TypeScript factory pattern
func StringProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Description: description,
	}
}

// RequiredStringProperty creates a required string property - TypeScript factory pattern
func RequiredStringProperty(description string) *jsonschema.Schema {
	minLen := 1
	return &jsonschema.Schema{
		Type:        "string",
		Description: description,
		MinLength:   &minLen, // Non-empty required
	}
}

// StringPropertyWithPattern creates a string property with pattern - TypeScript regex pattern
func StringPropertyWithPattern(description, pattern string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Description: description,
		Pattern:     pattern,
	}
}

// IntegerProperty creates an integer property schema - TypeScript number pattern
func IntegerProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "integer",
		Description: description,
		Minimum:     float64Ptr(0), // Default to non-negative
	}
}

// IntegerPropertyWithRange creates an integer property with range - TypeScript range pattern
func IntegerPropertyWithRange(description string, min, max *float64) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "integer",
		Description: description,
		Minimum:     min,
		Maximum:     max,
	}
}

// BooleanProperty creates a boolean property schema - TypeScript boolean pattern
func BooleanProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "boolean",
		Description: description,
	}
}

// ArrayProperty creates an array property schema - TypeScript array pattern
func ArrayProperty(description string, items *jsonschema.Schema) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "array",
		Description: description,
		Items:       items,
	}
}

// ArrayPropertyWithConstraints creates an array property with constraints - TypeScript array with constraints
func ArrayPropertyWithConstraints(description string, items *jsonschema.Schema, minItems, maxItems *int) *jsonschema.Schema {
	schema := &jsonschema.Schema{
		Type:        "array",
		Description: description,
		Items:       items,
	}

	if minItems != nil {
		schema.MinItems = minItems
	}
	if maxItems != nil {
		schema.MaxItems = maxItems
	}

	return schema
}

// ObjectProperty creates an object property schema - TypeScript object pattern
func ObjectProperty(description string, properties map[string]*jsonschema.Schema) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "object",
		Description: description,
		Properties:  properties,
	}
}

// EnumProperty creates an enum property schema - TypeScript enum pattern
func EnumProperty(description string, values []interface{}) *jsonschema.Schema {
	return &jsonschema.Schema{
		Description: description,
		Enum:        values,
	}
}

// OneOfProperty creates a oneOf property schema - TypeScript union type pattern
func OneOfProperty(description string, schemas []*jsonschema.Schema) *jsonschema.Schema {
	return &jsonschema.Schema{
		Description: description,
		OneOf:       schemas,
	}
}

// Helper function for float64 pointer - TypeScript utility pattern
func float64Ptr(f float64) *float64 {
	return &f
}
