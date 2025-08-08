package tools

import (
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/webex"
)

// CoreListRoomsParams for listing rooms - minimal version
type CoreListRoomsParams struct {
	Type string `json:"type,omitempty" query:"type"`
	Max  int    `json:"max,omitempty" query:"max" includeZero:"false"`
}

// NewListRoomsTool creates a tool to list rooms - essential for finding conversation spaces
func NewListRoomsTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"type": StringProperty("direct (1:1), group (group space)."),
		"max":  IntegerProperty("Limit the maximum number of rooms in the response."),
	}

	return NewListTool[CoreListRoomsParams](
		"list_rooms",
		"List rooms visible to the authenticated user.",
		"/rooms",
		properties,
		[]string{}, // No required fields for listing rooms
	)
}

// NewGetMyOwnDetailsTool gets current user details - essential for bot identity
func NewGetMyOwnDetailsTool() Tool {
	return NewGenericTool("get_my_own_details", "Get details for the authenticated user.",
		SimpleSchema("Get details for the authenticated user.", nil, nil),
		func(params *map[string]interface{}, client webex.HTTPClient) (interface{}, error) {
			return client.Get("/people/me", nil)
		})
}
