package tools

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestExecuteWithMapBase(t *testing.T) {
	// Create mock tool that implements ToolWithExecute
	mockTool := &mockToolWithExecute{
		executeFunc: func(args json.RawMessage) (interface{}, error) {
			// Parse args
			var argsMap map[string]interface{}
			if err := json.Unmarshal(args, &argsMap); err != nil {
				return nil, err
			}

			// Return test response
			return map[string]interface{}{
				"result": "success",
				"id":     "123",
				"args":   argsMap,
			}, nil
		},
	}

	// Test successful execution
	args := map[string]interface{}{
		"key": "value",
	}
	result, err := ExecuteWithMapBase(mockTool, args)
	if err != nil {
		t.Fatalf("ExecuteWithMapBase() error = %v", err)
	}

	// Check result
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["result"] != "success" {
		t.Errorf("Expected result = success, got %v", resultMap["result"])
	}

	// Verify args were passed correctly
	if argsInResult, ok := resultMap["args"].(map[string]interface{}); ok {
		if argsInResult["key"] != "value" {
			t.Errorf("Args not passed correctly: %v", argsInResult)
		}
	}
}

func TestExecuteWithMapBase_ErrorHandling(t *testing.T) {
	// Create mock tool that returns an error
	mockTool := &mockToolWithExecute{
		executeFunc: func(args json.RawMessage) (interface{}, error) {
			return nil, fmt.Errorf("execution failed")
		},
	}

	// Test error handling
	_, err := ExecuteWithMapBase(mockTool, map[string]interface{}{})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "execution failed" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestExecuteWithMapBase_MarshalError(t *testing.T) {
	// Create mock tool
	mockTool := &mockToolWithExecute{
		executeFunc: func(args json.RawMessage) (interface{}, error) {
			return nil, nil
		},
	}

	// Test with unmarshalable args (channel causes marshal error)
	ch := make(chan int)
	args := map[string]interface{}{
		"channel": ch,
	}

	_, err := ExecuteWithMapBase(mockTool, args)
	if err == nil {
		t.Fatal("Expected marshal error, got nil")
	}

	if !contains(err.Error(), "failed to marshal arguments") {
		t.Errorf("Unexpected error: %v", err)
	}
}

// mockToolWithExecute implements ToolWithExecute for testing
type mockToolWithExecute struct {
	executeFunc func(json.RawMessage) (interface{}, error)
}

func (m *mockToolWithExecute) Execute(args json.RawMessage) (interface{}, error) {
	if m.executeFunc != nil {
		return m.executeFunc(args)
	}
	return nil, nil
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
