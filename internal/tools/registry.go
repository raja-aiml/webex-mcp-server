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

// LoadCoreTools creates registry with only essential conversation tools
// Following KISS principle - only minimum required for bot conversations
func LoadCoreTools() (*Registry, error) {
	registry := NewRegistry()
	manager := NewPluginManager()

	// Load only core plugins for conversation functionality
	LoadCorePlugins(manager)

	// Load plugins into registry
	if err := manager.LoadPlugins(registry); err != nil {
		return nil, err
	}

	return registry, nil
}

// LoadAllTools creates registry with both core and advanced tools
// Used when full functionality is needed
func LoadAllTools() (*Registry, error) {
	registry := NewRegistry()
	manager := NewPluginManager()

	// Load core plugins first
	LoadCorePlugins(manager)

	// Load all advanced plugins
	LoadDefaultPlugins(manager)

	// Load plugins into registry
	if err := manager.LoadPlugins(registry); err != nil {
		return nil, err
	}

	return registry, nil
}

// LoadTools is deprecated, use LoadCoreTools or LoadAllTools instead
// Kept for backward compatibility
func LoadTools() (*Registry, error) {
	return LoadAllTools()
}
