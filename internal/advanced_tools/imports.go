package advanced_tools

import (
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

// Re-export base types for advanced tools to use
type Tool = tools.Tool
type Registry = tools.Registry

// Helper function aliases - these are non-generic functions that can be aliased
var (
	SimpleSchema    = tools.SimpleSchema
	StringProperty  = tools.StringProperty
	IntegerProperty = tools.IntegerProperty
	BooleanProperty = tools.BooleanProperty
	ArrayProperty   = tools.ArrayProperty
	ObjectProperty  = tools.ObjectProperty
)
