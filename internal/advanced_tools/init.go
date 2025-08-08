package advanced_tools

import (
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

func init() {
	// Register the advanced plugin loader with the tools package
	// This avoids circular dependency by registering at runtime
	tools.SetAdvancedLoader(LoadAdvancedPlugins)
}
