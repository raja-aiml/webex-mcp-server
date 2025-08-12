package tools

// LoadDefaultPlugins loads all advanced tool plugins
// This now uses the registered loader to avoid circular dependency
func LoadDefaultPlugins(manager *PluginManager) {
	// Load advanced plugins using the registered loader
	LoadAdvancedPlugins(manager)
}
