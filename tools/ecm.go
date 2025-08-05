package tools

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// ECM (Enterprise Content Management) Tools - CRUD operations for ECM folders

// CreateECMFolderConfigurationTool creates an ECM folder configuration
type CreateECMFolderConfigurationTool struct {
	ToolBase
}

func NewCreateECMFolderConfigurationTool() *CreateECMFolderConfigurationTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId":      StringProperty("A unique identifier for the room."),
		"folderId":    StringProperty("The ECM folder ID."),
		"displayName": StringProperty("A user-friendly name for the ECM folder."),
	}, []string{"roomId", "folderId"})

	return &CreateECMFolderConfigurationTool{
		ToolBase: NewToolBase("create_an_ecm_folder_configuration", "Create an ECM folder configuration", schema),
	}
}

func (t *CreateECMFolderConfigurationTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Validate required fields
	if params["roomId"] == nil || params["roomId"] == "" {
		return nil, fmt.Errorf("roomId is required")
	}
	if params["folderId"] == nil || params["folderId"] == "" {
		return nil, fmt.Errorf("folderId is required")
	}

	return t.client.Post("/rooms/linkedFolders", params)
}

func (t *CreateECMFolderConfigurationTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// GetECMFolderDetailsTool gets ECM folder details
type GetECMFolderDetailsTool struct {
	ToolBase
}

func NewGetECMFolderDetailsTool() *GetECMFolderDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"linkedFolderId": StringProperty("The unique identifier for the linked folder."),
	}, []string{"linkedFolderId"})

	return &GetECMFolderDetailsTool{
		ToolBase: NewToolBase("get_ecm_folder_details", "Get ECM folder details", schema),
	}
}

func (t *GetECMFolderDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		LinkedFolderId string `json:"linkedFolderId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.LinkedFolderId == "" {
		return nil, fmt.Errorf("linkedFolderId is required")
	}

	return t.client.Get(fmt.Sprintf("/rooms/linkedFolders/%s", params.LinkedFolderId), nil)
}

func (t *GetECMFolderDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// ListECMFolderTool lists ECM folders
type ListECMFolderTool struct {
	ToolBase
}

func NewListECMFolderTool() *ListECMFolderTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId": StringProperty("List linked folders associated with a room, by ID."),
	}, []string{"roomId"})

	return &ListECMFolderTool{
		ToolBase: NewToolBase("list_ecm_folder", "List ECM folders", schema),
	}
}

func (t *ListECMFolderTool) Execute(args json.RawMessage) (interface{}, error) {
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
	return t.client.Get("/rooms/linkedFolders", queryParams)
}

func (t *ListECMFolderTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// UpdateECMLinkedFolderTool updates an ECM linked folder
type UpdateECMLinkedFolderTool struct {
	ToolBase
}

func NewUpdateECMLinkedFolderTool() *UpdateECMLinkedFolderTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"linkedFolderId": StringProperty("The unique identifier for the linked folder."),
		"displayName":    StringProperty("A user-friendly name for the ECM folder."),
	}, []string{"linkedFolderId"})

	return &UpdateECMLinkedFolderTool{
		ToolBase: NewToolBase("update_an_ecm_linked_folder", "Update an ECM linked folder", schema),
	}
}

func (t *UpdateECMLinkedFolderTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	linkedFolderId, ok := params["linkedFolderId"].(string)
	if !ok || linkedFolderId == "" {
		return nil, fmt.Errorf("linkedFolderId is required")
	}

	// Remove linkedFolderId from params as it's in the URL
	delete(params, "linkedFolderId")

	// At least one field to update must be specified
	if len(params) == 0 {
		return nil, fmt.Errorf("at least one field to update must be specified")
	}

	return t.client.Put(fmt.Sprintf("/rooms/linkedFolders/%s", linkedFolderId), params)
}

func (t *UpdateECMLinkedFolderTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// UnlinkECMLinkedFolderTool unlinks an ECM linked folder
type UnlinkECMLinkedFolderTool struct {
	ToolBase
}

func NewUnlinkECMLinkedFolderTool() *UnlinkECMLinkedFolderTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"linkedFolderId": StringProperty("The unique identifier for the linked folder."),
	}, []string{"linkedFolderId"})

	return &UnlinkECMLinkedFolderTool{
		ToolBase: NewToolBase("unlink_an_ecm_linked_folder", "Unlink an ECM linked folder", schema),
	}
}

func (t *UnlinkECMLinkedFolderTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		LinkedFolderId string `json:"linkedFolderId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.LinkedFolderId == "" {
		return nil, fmt.Errorf("linkedFolderId is required")
	}

	if err := t.client.Delete(fmt.Sprintf("/rooms/linkedFolders/%s", params.LinkedFolderId)); err != nil {
		return nil, err
	}

	return map[string]interface{}{"success": true, "message": "ECM folder unlinked"}, nil
}

func (t *UnlinkECMLinkedFolderTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}