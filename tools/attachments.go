package tools

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// Attachment Action Tools - Operations for attachment actions

// CreateAttachmentActionTool creates an attachment action
type CreateAttachmentActionTool struct {
	ToolBase
}

func NewCreateAttachmentActionTool() *CreateAttachmentActionTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"type":      StringProperty("The type of action."),
		"messageId": StringProperty("The ID of the message with attachment."),
		"inputs":    ObjectProperty("The attachment action's inputs."),
	}, []string{"type", "messageId"})

	return &CreateAttachmentActionTool{
		ToolBase: NewToolBase("create_an_attachment_action", "Create an attachment action", schema),
	}
}

func (t *CreateAttachmentActionTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Validate required fields
	if params["type"] == nil || params["type"] == "" {
		return nil, fmt.Errorf("type is required")
	}
	if params["messageId"] == nil || params["messageId"] == "" {
		return nil, fmt.Errorf("messageId is required")
	}

	return t.client.Post("/attachment/actions", params)
}

func (t *CreateAttachmentActionTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// GetAttachmentActionDetailsTool gets attachment action details
type GetAttachmentActionDetailsTool struct {
	ToolBase
}

func NewGetAttachmentActionDetailsTool() *GetAttachmentActionDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"actionId": StringProperty("A unique identifier for the attachment action."),
	}, []string{"actionId"})

	return &GetAttachmentActionDetailsTool{
		ToolBase: NewToolBase("get_attachment_action_details", "Get attachment action details", schema),
	}
}

func (t *GetAttachmentActionDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		ActionId string `json:"actionId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.ActionId == "" {
		return nil, fmt.Errorf("actionId is required")
	}

	return t.client.Get(fmt.Sprintf("/attachment/actions/%s", params.ActionId), nil)
}

func (t *GetAttachmentActionDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}