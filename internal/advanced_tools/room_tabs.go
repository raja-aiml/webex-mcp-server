package advanced_tools

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

type ListRoomTabsParams struct {
	RoomId string `json:"roomId" required:"true"`
}

// CreateRoomTabParams defines the parameters for creating a room tab
type CreateRoomTabParams struct {
	RoomId      string `json:"roomId" required:"true"`
	ContentUrl  string `json:"contentUrl" required:"true"`
	DisplayName string `json:"displayName" required:"true"`
}

// UpdateRoomTabParams defines the parameters for updating a room tab
type UpdateRoomTabParams struct {
	RoomTabId   string `json:"roomTabId" required:"true"`
	ContentUrl  string `json:"contentUrl,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// NewListRoomTabsTool lists room tabs
func NewListRoomTabsTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"roomId": StringProperty("List tabs for a room, by room ID."),
	}

	return tools.NewListTool[ListRoomTabsParams](
		"list_room_tabs",
		"List tabs for a room.",
		"/roomTabs",
		properties,
		[]string{"roomId"}, // roomId is required
	)
}

// NewCreateRoomTabTool creates a new room tab
func NewCreateRoomTabTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"roomId":      StringProperty("The room ID."),
		"contentUrl":  StringProperty("URL of the tab content."),
		"displayName": StringProperty("User-friendly name for the tab."),
	}

	return tools.NewCreateTool[CreateRoomTabParams](
		"create_a_room_tab",
		"Add a tab to a room.",
		"/roomTabs",
		properties,
		[]string{"roomId", "contentUrl", "displayName"},
	)
}

// NewGetRoomTabDetailsTool gets room tab details
func NewGetRoomTabDetailsTool() Tool {
	return tools.NewGetTool(
		"get_room_tab_details",
		"Get details for a room tab by ID.",
		"/roomTabs",
		"roomTabId",
		"The unique identifier for the room tab.",
	)
}

// NewUpdateRoomTabTool updates a room tab
func NewUpdateRoomTabTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"roomTabId":   StringProperty("The unique identifier for the room tab."),
		"contentUrl":  StringProperty("URL of the tab content."),
		"displayName": StringProperty("User-friendly name for the tab."),
	}

	return tools.NewUpdateTool[UpdateRoomTabParams](
		"update_a_room_tab",
		"Update a room tab by ID.",
		"/roomTabs",
		"roomTabId",
		properties,
		[]string{"roomTabId"},
	)
}

// NewDeleteRoomTabTool deletes a room tab
func NewDeleteRoomTabTool() Tool {
	return tools.NewDeleteTool(
		"delete_a_room_tab",
		"Delete a room tab by ID.",
		"/roomTabs",
		"roomTabId",
		"The unique identifier for the room tab.",
	)
}
