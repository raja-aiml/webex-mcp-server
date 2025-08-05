package tools

// ToolPlugin defines the interface for tool plugins
// This implements Open/Closed Principle - open for extension, closed for modification
type ToolPlugin interface {
	// Register adds all tools provided by this plugin to the registry
	Register(registry *Registry) error
	// Name returns the plugin name
	Name() string
	// Version returns the plugin version
	Version() string
}

// PluginManager manages tool plugins
type PluginManager struct {
	plugins []ToolPlugin
}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make([]ToolPlugin, 0),
	}
}

// RegisterPlugin adds a plugin to the manager
func (pm *PluginManager) RegisterPlugin(plugin ToolPlugin) {
	pm.plugins = append(pm.plugins, plugin)
}

// LoadPlugins loads all registered plugins into the registry
func (pm *PluginManager) LoadPlugins(registry *Registry) error {
	for _, plugin := range pm.plugins {
		if err := plugin.Register(registry); err != nil {
			return err
		}
	}
	return nil
}

// GetPlugins returns all registered plugins
func (pm *PluginManager) GetPlugins() []ToolPlugin {
	return pm.plugins
}
