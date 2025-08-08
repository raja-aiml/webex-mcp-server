package advanced_tools

import (
	"fmt"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
	"github.com/raja-aiml/webex-mcp-server/internal/webex"
)

type ListRoomsParams struct {
	TeamId string `json:"teamId,omitempty" query:"teamId"`
	Type   string `json:"type,omitempty" query:"type"`
	SortBy string `json:"sortBy,omitempty" query:"sortBy"`
	Max    int    `json:"max,omitempty" query:"max" includeZero:"false"`
}

// CreateRoomParams defines the parameters for creating a room
type CreateRoomParams struct {
	Title              string `json:"title" required:"true"`
	TeamId             string `json:"teamId,omitempty"`
	ClassificationId   string `json:"classificationId,omitempty"`
	IsLocked           bool   `json:"isLocked,omitempty"`
	IsPublic           bool   `json:"isPublic,omitempty"`
	Description        string `json:"description,omitempty"`
	IsAnnouncementOnly bool   `json:"isAnnouncementOnly,omitempty"`
}

// UpdateRoomParams defines the parameters for updating a room
type UpdateRoomParams struct {
	RoomId             string `json:"roomId" required:"true"`
	Title              string `json:"title,omitempty"`
	ClassificationId   string `json:"classificationId,omitempty"`
	TeamId             string `json:"teamId,omitempty"`
	IsLocked           bool   `json:"isLocked,omitempty"`
	IsPublic           bool   `json:"isPublic,omitempty"`
	Description        string `json:"description,omitempty"`
	IsAnnouncementOnly bool   `json:"isAnnouncementOnly,omitempty"`
	IsReadOnly         bool   `json:"isReadOnly,omitempty"`
}

// NewCreateRoomTool creates a new Webex room
func NewCreateRoomTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"title":              StringProperty("A user-friendly name for the room."),
		"teamId":             StringProperty("The ID for the team with which this room is associated."),
		"classificationId":   StringProperty("The classification ID for the room."),
		"isLocked":           BooleanProperty("Whether the room is locked (moderator approval required)."),
		"isPublic":           BooleanProperty("Whether the room is public (allows guest users)."),
		"description":        StringProperty("The description of the room."),
		"isAnnouncementOnly": BooleanProperty("Whether the room is announcement only."),
	}

	return tools.NewCreateTool[CreateRoomParams](
		"create_a_room",
		"Create a new Webex room.",
		"/rooms",
		properties,
		[]string{"title"},
	)
}

// NewGetRoomDetailsTool gets details of a specific room
func NewGetRoomDetailsTool() Tool {
	return tools.NewGetTool(
		"get_room_details",
		"Get details of a specific room.",
		"/rooms",
		"roomId",
		"The unique identifier for the room.",
	)
}

// NewUpdateRoomTool updates a room
func NewUpdateRoomTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"roomId":             StringProperty("The unique identifier for the room."),
		"title":              StringProperty("A user-friendly name for the room."),
		"classificationId":   StringProperty("The classification ID for the room."),
		"teamId":             StringProperty("The teamId to which this room belongs."),
		"isLocked":           BooleanProperty("Whether the room is locked (moderator approval required)."),
		"isPublic":           BooleanProperty("Whether the room is public (allows guest users)."),
		"description":        StringProperty("The description of the room."),
		"isAnnouncementOnly": BooleanProperty("Whether the room is announcement only."),
		"isReadOnly":         BooleanProperty("Whether the room is read only."),
	}

	return tools.NewUpdateTool[UpdateRoomParams](
		"update_a_room",
		"Update a room.",
		"/rooms",
		"roomId",
		properties,
		[]string{"roomId"},
	)
}

// NewDeleteRoomTool deletes a room
func NewDeleteRoomTool() Tool {
	return tools.NewDeleteTool(
		"delete_a_room",
		"Delete a room.",
		"/rooms",
		"roomId",
		"The unique identifier for the room.",
	)
}

// NewGetRoomMeetingDetailsTool gets meeting details for a room
func NewGetRoomMeetingDetailsTool() Tool {
	schema := SimpleSchema("Get meeting details for a room.", map[string]*jsonschema.Schema{
		"roomId": StringProperty("The unique identifier for the room."),
	}, []string{"roomId"})

	return tools.NewGenericTool("get_room_meeting_details", "Get meeting details for a room.", schema,
		func(params *map[string]interface{}, client webex.HTTPClient) (interface{}, error) {
			roomId, ok := (*params)["roomId"].(string)
			if !ok || roomId == "" {
				return nil, fmt.Errorf("roomId is required")
			}
			return client.Get(fmt.Sprintf("/rooms/%s/meetingInfo", roomId), nil)
		})
}
