package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// ListRoomsTool lists Webex rooms
type ListRoomsTool struct {
	ToolBase
}

func NewListRoomsTool() *ListRoomsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"teamId":      StringProperty("List rooms associated with a team, by ID."),
		"type":        StringProperty("List rooms by type: 'direct' or 'group'."),
		"sortBy":      StringProperty("Sort results by: 'id', 'lastactivity', or 'created'."),
		"max":         IntegerProperty("Limit the maximum number of rooms in the response."),
	}, []string{})

	return &ListRoomsTool{
		ToolBase: NewToolBase("list_rooms", "List Webex rooms.", schema),
	}
}

func (t *ListRoomsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		TeamId  string `json:"teamId,omitempty"`
		Type    string `json:"type,omitempty"`
		SortBy  string `json:"sortBy,omitempty"`
		Max     int    `json:"max,omitempty"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	queryParams := map[string]string{}
	if params.TeamId != "" {
		queryParams["teamId"] = params.TeamId
	}
	if params.Type != "" {
		queryParams["type"] = params.Type
	}
	if params.SortBy != "" {
		queryParams["sortBy"] = params.SortBy
	}
	if params.Max > 0 {
		queryParams["max"] = strconv.Itoa(params.Max)
	}

	return t.client.Get("/rooms", queryParams)
}

func (t *ListRoomsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// CreateRoomTool creates a new Webex room
type CreateRoomTool struct {
	ToolBase
}

func NewCreateRoomTool() *CreateRoomTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"title":               StringProperty("A user-friendly name for the room."),
		"teamId":              StringProperty("The ID for the team with which this room is associated."),
		"classificationId":    StringProperty("The classification ID for the room."),
		"isLocked":            BooleanProperty("Whether the room is locked (moderator approval required)."),
		"isPublic":            BooleanProperty("Whether the room is public (allows guest users)."),
		"description":         StringProperty("The description of the room."),
		"isAnnouncementOnly":  BooleanProperty("Whether the room is announcement only."),
	}, []string{"title"})

	return &CreateRoomTool{
		ToolBase: NewToolBase("create_a_room", "Create a new Webex room.", schema),
	}
}

func (t *CreateRoomTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params["title"] == nil || params["title"] == "" {
		return nil, fmt.Errorf("title is required")
	}

	return t.client.Post("/rooms", params)
}

func (t *CreateRoomTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// GetRoomDetailsTool gets details of a specific room
type GetRoomDetailsTool struct {
	ToolBase
}

func NewGetRoomDetailsTool() *GetRoomDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId": StringProperty("The unique identifier for the room."),
	}, []string{"roomId"})

	return &GetRoomDetailsTool{
		ToolBase: NewToolBase("get_room_details", "Get details of a specific room.", schema),
	}
}

func (t *GetRoomDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		RoomId string `json:"roomId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.RoomId == "" {
		return nil, fmt.Errorf("roomId is required")
	}

	return t.client.Get(fmt.Sprintf("/rooms/%s", params.RoomId), nil)
}

func (t *GetRoomDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// UpdateRoomTool updates a room
type UpdateRoomTool struct {
	ToolBase
}

func NewUpdateRoomTool() *UpdateRoomTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId":              StringProperty("The unique identifier for the room."),
		"title":               StringProperty("A user-friendly name for the room."),
		"classificationId":    StringProperty("The classification ID for the room."),
		"teamId":              StringProperty("The teamId to which this room belongs."),
		"isLocked":            BooleanProperty("Whether the room is locked (moderator approval required)."),
		"isPublic":            BooleanProperty("Whether the room is public (allows guest users)."),
		"description":         StringProperty("The description of the room."),
		"isAnnouncementOnly":  BooleanProperty("Whether the room is announcement only."),
		"isReadOnly":          BooleanProperty("Whether the room is read only."),
	}, []string{"roomId"})

	return &UpdateRoomTool{
		ToolBase: NewToolBase("update_a_room", "Update a room.", schema),
	}
}

func (t *UpdateRoomTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	roomId, ok := params["roomId"].(string)
	if !ok || roomId == "" {
		return nil, fmt.Errorf("roomId is required")
	}

	// Remove roomId from params as it's in the URL
	delete(params, "roomId")

	// At least one field to update must be specified
	if len(params) == 0 {
		return nil, fmt.Errorf("at least one field to update must be specified")
	}

	return t.client.Put(fmt.Sprintf("/rooms/%s", roomId), params)
}

func (t *UpdateRoomTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// DeleteRoomTool deletes a room
type DeleteRoomTool struct {
	ToolBase
}

func NewDeleteRoomTool() *DeleteRoomTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId": StringProperty("The unique identifier for the room."),
	}, []string{"roomId"})

	return &DeleteRoomTool{
		ToolBase: NewToolBase("delete_a_room", "Delete a room.", schema),
	}
}

func (t *DeleteRoomTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		RoomId string `json:"roomId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.RoomId == "" {
		return nil, fmt.Errorf("roomId is required")
	}

	err := t.client.Delete(fmt.Sprintf("/rooms/%s", params.RoomId))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"success": true}, nil
}

func (t *DeleteRoomTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// GetRoomMeetingDetailsTool gets meeting details for a room
type GetRoomMeetingDetailsTool struct {
	ToolBase
}

func NewGetRoomMeetingDetailsTool() *GetRoomMeetingDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId": StringProperty("The unique identifier for the room."),
	}, []string{"roomId"})

	return &GetRoomMeetingDetailsTool{
		ToolBase: NewToolBase("get_room_meeting_details", "Get meeting details for a room.", schema),
	}
}

func (t *GetRoomMeetingDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		RoomId string `json:"roomId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.RoomId == "" {
		return nil, fmt.Errorf("roomId is required")
	}

	return t.client.Get(fmt.Sprintf("/rooms/%s/meetingInfo", params.RoomId), nil)
}

func (t *GetRoomMeetingDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}