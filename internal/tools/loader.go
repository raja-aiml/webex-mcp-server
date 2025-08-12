package tools

// LoaderFunc is a function type for loading plugins
type LoaderFunc func(*PluginManager)

// advancedLoader holds the function for loading advanced plugins
var advancedLoader LoaderFunc

// SetAdvancedLoader sets the loader function for advanced plugins
// This allows advanced_tools package to register its loader without circular dependency
func SetAdvancedLoader(loader LoaderFunc) {
	advancedLoader = loader
}

// LoadAdvancedPlugins loads advanced plugins if loader is set
func LoadAdvancedPlugins(manager *PluginManager) {
	if advancedLoader != nil {
		advancedLoader(manager)
	}
}
