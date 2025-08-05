package tools

import (
	"encoding/json"
	"fmt"
)

// ToolWithExecute interface for tools that have Execute method
type ToolWithExecute interface {
	Execute(args json.RawMessage) (interface{}, error)
}

// ExecuteWithMapBase provides the common ExecuteWithMap implementation
// This function should be called by concrete tool types that embed ToolBase
func ExecuteWithMapBase(tool ToolWithExecute, args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return tool.Execute(argsJSON)
}
