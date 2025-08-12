package tools

import (
	"testing"

	"github.com/raja-aiml/webex-mcp-server/internal/config"
	"github.com/raja-aiml/webex-mcp-server/internal/testutil"
)

func TestNewCreateMessageTool(t *testing.T) {
	// Set up environment for default client
	config.ResetForTesting()
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-key")
	defer func() {
		cleanup()
		config.ResetForTesting()
		defaultClient = nil
	}()

	// Initialize default client
	_ = InitializeDefaultClient()

	tool := NewCreateMessageTool()

	if tool == nil {
		t.Fatal("NewCreateMessageTool returned nil")
	}

	// Tool names are now snake_case based on implementation
	if tool.Name() != "create_a_message" {
		t.Errorf("Name() = %v, want create_a_message", tool.Name())
	}
}

func TestNewListDirectMessagesTool(t *testing.T) {
	// Set up environment for default client
	config.ResetForTesting()
	cleanup := testutil.SetEnv(t, "WEBEX_PUBLIC_WORKSPACE_API_KEY", "test-key")
	defer func() {
		cleanup()
		config.ResetForTesting()
		defaultClient = nil
	}()

	// Initialize default client
	_ = InitializeDefaultClient()

	tool := NewListDirectMessagesTool()

	if tool == nil {
		t.Fatal("NewListDirectMessagesTool returned nil")
	}

	if tool.Name() != "list_direct_messages" {
		t.Errorf("Name() = %v, want list_direct_messages", tool.Name())
	}
}

// TestQueryParamsFunction removed - QueryParams now only works with maps, not structs
/*
func TestQueryParamsFunction(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		want   map[string]string
		wantNil bool
	}{
		{
			name: "struct with query tags",
			input: struct {
				RoomId string `query:"roomId"`
				Max    int    `query:"max"`
				Skip   int    // no tag
			}{
				RoomId: "test-room",
				Max:    10,
				Skip:   5,
			},
			want: map[string]string{
				"roomId": "test-room",
				"max":    "10",
				"Skip":   "5", // Without tag, field name is used
			},
		},
		{
			name: "struct with empty values",
			input: struct {
				RoomId string `query:"roomId"`
				Max    int    `query:"max"`
			}{
				RoomId: "",
				Max:    0,
			},
			want: map[string]string{},
		},
		{
			name:    "non-struct input",
			input:   "not a struct",
			wantNil: true,
		},
		{
			name:    "nil input",
			input:   nil,
			wantNil: true,
		},
		{
			name: "pointer to struct",
			input: &struct {
				ID string `query:"id"`
			}{
				ID: "test-id",
			},
			want: map[string]string{
				"id": "test-id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QueryParams(tt.input)

			if tt.wantNil {
				if got != nil {
					t.Errorf("QueryParams() = %v, want nil", got)
				}
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("QueryParams() returned %d params, want %d", len(got), len(tt.want))
			}

			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("QueryParams()[%s] = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}
*/
