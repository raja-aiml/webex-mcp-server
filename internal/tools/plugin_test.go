package tools

import (
	"encoding/json"
	"testing"
)

func TestNewPluginManager(t *testing.T) {
	pm := NewPluginManager()

	if pm == nil {
		t.Fatal("Expected non-nil PluginManager")
	}

	// Check that plugins map is initialized
	plugins := pm.GetPlugins()
	if plugins == nil {
		t.Error("Expected non-nil plugins slice")
	}

	if len(plugins) != 0 {
		t.Errorf("Expected empty plugins slice, got %d plugins", len(plugins))
	}
}

func TestPluginManager_RegisterPlugin(t *testing.T) {
	pm := NewPluginManager()

	// Create a test plugin
	plugin := &testPlugin{
		name:    "test-plugin",
		version: "1.0.0",
	}

	pm.RegisterPlugin(plugin)

	plugins := pm.GetPlugins()
	if len(plugins) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(plugins))
	}

	if len(plugins) > 0 && plugins[0] != plugin {
		t.Error("Registered plugin does not match")
	}
}

func TestPluginManager_LoadPlugins(t *testing.T) {
	pm := NewPluginManager()

	// Create test plugins
	plugin1 := &testPlugin{
		name:    "plugin1",
		version: "1.0.0",
		tools: []Tool{
			&testTool{name: "tool1", description: "Tool 1"},
			&testTool{name: "tool2", description: "Tool 2"},
		},
	}

	plugin2 := &testPlugin{
		name:    "plugin2",
		version: "2.0.0",
		tools: []Tool{
			&testTool{name: "tool3", description: "Tool 3"},
		},
	}

	pm.RegisterPlugin(plugin1)
	pm.RegisterPlugin(plugin2)

	// Load plugins into registry
	registry := NewRegistry()
	pm.LoadPlugins(registry)

	// Verify all tools are registered
	allTools := registry.GetTools()
	if len(allTools) != 3 {
		t.Errorf("Expected 3 tools, got %d", len(allTools))
	}

	// Verify individual tools
	tool1, exists := registry.GetTool("tool1")
	if !exists {
		t.Error("tool1 not found in registry")
	} else if tool1.Name() != "tool1" {
		t.Errorf("Expected tool name 'tool1', got %s", tool1.Name())
	}

	_, exists = registry.GetTool("tool3")
	if !exists {
		t.Error("tool3 not found in registry")
	}
}

func TestPluginManager_MultipleRegistrations(t *testing.T) {
	pm := NewPluginManager()

	plugin1 := &testPlugin{name: "plugin", version: "1.0.0"}
	plugin2 := &testPlugin{name: "plugin", version: "2.0.0"}

	pm.RegisterPlugin(plugin1)
	pm.RegisterPlugin(plugin2) // Should overwrite

	plugins := pm.GetPlugins()
	if len(plugins) != 2 {
		t.Errorf("Expected 2 plugins (both are added), got %d", len(plugins))
	}

	// Both plugins should be present - no overwriting in slice
	if plugins[0].Version() != "1.0.0" || plugins[1].Version() != "2.0.0" {
		t.Error("Expected both plugins to be present")
	}
}

// Test plugin implementation
type testPlugin struct {
	name    string
	version string
	tools   []Tool
}

func (p *testPlugin) Name() string    { return p.name }
func (p *testPlugin) Version() string { return p.version }
func (p *testPlugin) Register(registry *Registry) error {
	for _, tool := range p.tools {
		registry.Register(tool)
	}
	return nil
}

// Test tool implementation
type testTool struct {
	name        string
	description string
}

func (t *testTool) Name() string                                      { return t.name }
func (t *testTool) Description() string                               { return t.description }
func (t *testTool) GetInputSchema() interface{}                       { return nil }
func (t *testTool) Execute(args json.RawMessage) (interface{}, error) { return nil, nil }
func (t *testTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return t.Execute(nil)
}
