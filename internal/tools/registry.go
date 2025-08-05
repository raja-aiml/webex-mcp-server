package tools

import (
	"encoding/json"
	"fmt"
)

type Tool interface {
	Name() string
	Description() string
	GetInputSchema() interface{}
	Execute(args json.RawMessage) (interface{}, error)
	ExecuteWithMap(args map[string]interface{}) (interface{}, error)
}

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Register(tool Tool) error {
	if _, exists := r.tools[tool.Name()]; exists {
		return fmt.Errorf("tool %s already registered", tool.Name())
	}
	r.tools[tool.Name()] = tool
	return nil
}

func (r *Registry) GetTool(name string) (Tool, bool) {
	tool, exists := r.tools[name]
	return tool, exists
}

func (r *Registry) GetTools() []Tool {
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// LoadTools creates and populates the tool registry
// The actual plugin loading is done by LoadDefaultPlugins function
func LoadTools() (*Registry, error) {
	registry := NewRegistry()

	// Use plugin architecture for extensibility
	// This implements Open/Closed Principle - open for extension, closed for modification
	manager := NewPluginManager()

	// Load all available plugins
	// This is defined in the same package to avoid circular imports
	LoadDefaultPlugins(manager)

	// Load plugins into registry
	if err := manager.LoadPlugins(registry); err != nil {
		return nil, err
	}

	return registry, nil
}
