package tools

import (
	"testing"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	if registry == nil {
		t.Fatal("Expected non-nil Registry")
	}

	// Check that tools slice is initialized
	tools := registry.GetTools()
	if tools == nil {
		t.Error("Expected non-nil tools slice")
	}

	if len(tools) != 0 {
		t.Errorf("Expected empty tools slice, got %d tools", len(tools))
	}
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()

	tool := &testTool{
		name:        "test-tool",
		description: "Test tool",
	}

	err := registry.Register(tool)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Verify tool is registered
	registeredTool, exists := registry.GetTool("test-tool")
	if !exists {
		t.Error("Tool not found after registration")
	} else if registeredTool != tool {
		t.Error("Registered tool does not match")
	}
}

func TestRegistry_Register_Duplicate(t *testing.T) {
	registry := NewRegistry()

	tool1 := &testTool{name: "duplicate", description: "First"}
	tool2 := &testTool{name: "duplicate", description: "Second"}

	// First registration should succeed
	err := registry.Register(tool1)
	if err != nil {
		t.Fatalf("First Register() error = %v", err)
	}

	// Second registration should fail (duplicate)
	err = registry.Register(tool2)
	if err == nil {
		t.Error("Expected error for duplicate registration")
	}

	// Verify first tool is still registered
	registeredTool, exists := registry.GetTool("duplicate")
	if !exists {
		t.Error("Tool not found after duplicate registration attempt")
	} else if registeredTool.Description() != "First" {
		t.Error("Expected first tool to remain registered")
	}
}

func TestRegistry_GetTool(t *testing.T) {
	registry := NewRegistry()

	// Register multiple tools
	tools := []Tool{
		&testTool{name: "tool1", description: "Tool 1"},
		&testTool{name: "tool2", description: "Tool 2"},
		&testTool{name: "tool3", description: "Tool 3"},
	}

	for _, tool := range tools {
		registry.Register(tool)
	}

	// Test retrieving existing tools
	for _, expectedTool := range tools {
		tool, exists := registry.GetTool(expectedTool.Name())
		if !exists {
			t.Errorf("GetTool(%s) returned not found", expectedTool.Name())
		} else if tool.Name() != expectedTool.Name() {
			t.Errorf("GetTool(%s) returned wrong tool: %s", expectedTool.Name(), tool.Name())
		}
	}

	// Test retrieving non-existent tool
	_, exists := registry.GetTool("non-existent")
	if exists {
		t.Error("GetTool() should return false for non-existent tool")
	}
}

func TestRegistry_GetTools(t *testing.T) {
	registry := NewRegistry()

	// Start with empty registry
	tools := registry.GetTools()
	if len(tools) != 0 {
		t.Errorf("Expected 0 tools in empty registry, got %d", len(tools))
	}

	// Register tools
	expectedTools := []Tool{
		&testTool{name: "tool1", description: "Tool 1"},
		&testTool{name: "tool2", description: "Tool 2"},
		&testTool{name: "tool3", description: "Tool 3"},
	}

	for _, tool := range expectedTools {
		registry.Register(tool)
	}

	// Get all tools
	allTools := registry.GetTools()
	if len(allTools) != len(expectedTools) {
		t.Errorf("Expected %d tools, got %d", len(expectedTools), len(allTools))
	}

	// Verify each tool is present
	foundTools := make(map[string]bool)
	for _, tool := range allTools {
		foundTools[tool.Name()] = true
	}

	for _, expected := range expectedTools {
		if !foundTools[expected.Name()] {
			t.Errorf("Tool %s not found in GetTools() result", expected.Name())
		}
	}
}

func TestLoadTools(t *testing.T) {
	// Test the LoadTools function that creates registry and loads plugins
	registry, err := LoadTools()

	if err != nil {
		t.Fatalf("LoadTools() error = %v", err)
	}

	if registry == nil {
		t.Fatal("LoadTools() returned nil registry")
	}

	// Registry should have tools loaded from default plugins
	tools := registry.GetTools()
	if len(tools) == 0 {
		t.Log("Registry has no tools loaded (plugins may not be configured)")
	}
}

func TestRegistry_CaseSensitivity(t *testing.T) {
	registry := NewRegistry()

	tool := &testTool{name: "TestTool", description: "Mixed case tool"}
	registry.Register(tool)

	// Exact case should work
	_, exists := registry.GetTool("TestTool")
	if !exists {
		t.Error("GetTool with exact case failed")
	}

	// Different case should not work (case-sensitive)
	_, exists = registry.GetTool("testtool")
	if exists {
		t.Error("GetTool should be case-sensitive")
	}

	_, exists = registry.GetTool("TESTTOOL")
	if exists {
		t.Error("GetTool should be case-sensitive")
	}
}

// testTool is already defined in plugin_test.go, so we'll reuse it
