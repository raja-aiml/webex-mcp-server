package server

import "fmt"

const (
	// MCPProtocolVersion is the version of the MCP protocol this server implements
	MCPProtocolVersion = "2024-11-05"

	// MaxToolNameLength is the maximum length for tool names
	MaxToolNameLength = 64

	// MaxDescriptionLength is the maximum length for descriptions
	MaxDescriptionLength = 1024
)

// ValidateToolName checks if a tool name is MCP-compliant
func ValidateToolName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("tool name cannot be empty")
	}
	if len(name) > MaxToolNameLength {
		return fmt.Errorf("tool name exceeds maximum length of %d characters", MaxToolNameLength)
	}
	// MCP recommends snake_case for tool names
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
			return fmt.Errorf("tool name must be snake_case (lowercase letters, numbers, and underscores only)")
		}
	}
	return nil
}

// ValidateToolDescription checks if a tool description is MCP-compliant
func ValidateToolDescription(description string) error {
	if len(description) == 0 {
		return fmt.Errorf("tool description cannot be empty")
	}
	if len(description) > MaxDescriptionLength {
		return fmt.Errorf("tool description exceeds maximum length of %d characters", MaxDescriptionLength)
	}
	return nil
}
