package tools

import (
	"encoding/json"
)

// BaseTool provides common functionality for all tools
type BaseTool struct{}

// ExecuteWithMap provides a default implementation that converts map to JSON
func (b *BaseTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	// This will be overridden by the actual tool implementation
	return nil, nil
}

// Helper method for tools to call from their ExecuteWithMap
func ExecuteWithMapHelper(tool interface{ Execute(json.RawMessage) (interface{}, error) }, args map[string]interface{}) (interface{}, error) {
	jsonData, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	return tool.Execute(jsonData)
}