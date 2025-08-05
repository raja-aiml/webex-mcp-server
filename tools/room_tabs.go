package tools

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// Room Tab Tools - CRUD operations for room tabs

// ListRoomTabsTool lists room tabs
type ListRoomTabsTool struct {
	ToolBase
}

func NewListRoomTabsTool() *ListRoomTabsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId": StringProperty("ID of the room."),
	}, []string{"roomId"})

	return &ListRoomTabsTool{
		ToolBase: NewToolBase("list_room_tabs", "List room tabs", schema),
	}
}

func (t *ListRoomTabsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		RoomId string `json:"roomId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.RoomId == "" {
		return nil, fmt.Errorf("roomId is required")
	}

	queryParams := map[string]string{"roomId": params.RoomId}
	return t.client.Get("/roomTabs", queryParams)
}

func (t *ListRoomTabsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// CreateRoomTabTool creates a room tab
type CreateRoomTabTool struct {
	ToolBase
}

func NewCreateRoomTabTool() *CreateRoomTabTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId":      StringProperty("A unique identifier for the room."),
		"contentUrl":  StringProperty("URL of the room tab content."),
		"displayName": StringProperty("User-friendly name for the room tab."),
	}, []string{"roomId", "contentUrl", "displayName"})

	return &CreateRoomTabTool{
		ToolBase: NewToolBase("create_a_room_tab", "Create a room tab", schema),
	}
}

func (t *CreateRoomTabTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Validate required fields
	required := []string{"roomId", "contentUrl", "displayName"}
	for _, field := range required {
		if params[field] == nil || params[field] == "" {
			return nil, fmt.Errorf("%s is required", field)
		}
	}

	return t.client.Post("/roomTabs", params)
}

func (t *CreateRoomTabTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// GetRoomTabDetailsTool gets room tab details
type GetRoomTabDetailsTool struct {
	ToolBase
}

func NewGetRoomTabDetailsTool() *GetRoomTabDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomTabId": StringProperty("The unique identifier for the room tab."),
	}, []string{"roomTabId"})

	return &GetRoomTabDetailsTool{
		ToolBase: NewToolBase("get_room_tab_details", "Get room tab details", schema),
	}
}

func (t *GetRoomTabDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		RoomTabId string `json:"roomTabId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.RoomTabId == "" {
		return nil, fmt.Errorf("roomTabId is required")
	}

	return t.client.Get(fmt.Sprintf("/roomTabs/%s", params.RoomTabId), nil)
}

func (t *GetRoomTabDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// UpdateRoomTabTool updates a room tab
type UpdateRoomTabTool struct {
	ToolBase
}

func NewUpdateRoomTabTool() *UpdateRoomTabTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomTabId":   StringProperty("The unique identifier for the room tab."),
		"contentUrl":  StringProperty("URL of the room tab content."),
		"displayName": StringProperty("User-friendly name for the room tab."),
	}, []string{"roomTabId"})

	return &UpdateRoomTabTool{
		ToolBase: NewToolBase("update_a_room_tab", "Update a room tab", schema),
	}
}

func (t *UpdateRoomTabTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	roomTabId, ok := params["roomTabId"].(string)
	if !ok || roomTabId == "" {
		return nil, fmt.Errorf("roomTabId is required")
	}

	// Remove roomTabId from params as it's in the URL
	delete(params, "roomTabId")

	// At least one field to update must be specified
	if len(params) == 0 {
		return nil, fmt.Errorf("at least one field to update must be specified")
	}

	return t.client.Put(fmt.Sprintf("/roomTabs/%s", roomTabId), params)
}

func (t *UpdateRoomTabTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// DeleteRoomTabTool deletes a room tab
type DeleteRoomTabTool struct {
	ToolBase
}

func NewDeleteRoomTabTool() *DeleteRoomTabTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomTabId": StringProperty("The unique identifier for the room tab."),
	}, []string{"roomTabId"})

	return &DeleteRoomTabTool{
		ToolBase: NewToolBase("delete_a_room_tab", "Delete a room tab", schema),
	}
}

func (t *DeleteRoomTabTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		RoomTabId string `json:"roomTabId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.RoomTabId == "" {
		return nil, fmt.Errorf("roomTabId is required")
	}

	if err := t.client.Delete(fmt.Sprintf("/roomTabs/%s", params.RoomTabId)); err != nil {
		return nil, err
	}

	return map[string]interface{}{"success": true, "message": "Room tab deleted"}, nil
}

func (t *DeleteRoomTabTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}