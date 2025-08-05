package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// ListMessagesTool lists messages in a Webex room
type ListMessagesTool struct {
	ToolBase
}

func NewListMessagesTool() *ListMessagesTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId":          StringProperty("The ID of the room to list messages from."),
		"parentId":        StringProperty("The ID of the parent message to filter by."),
		"mentionedPeople": StringProperty("List messages with these people mentioned, by ID."),
		"before":          StringProperty("List messages sent before a specific date and time."),
		"beforeMessage":   StringProperty("List messages sent before a specific message, by ID."),
		"max":             IntegerProperty("Limit the maximum number of messages in the response."),
	}, []string{"roomId"})

	return &ListMessagesTool{
		ToolBase: NewToolBase("list_messages", "List messages in a Webex room.", schema),
	}
}

func (t *ListMessagesTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		RoomId          string `json:"roomId"`
		ParentId        string `json:"parentId,omitempty"`
		MentionedPeople string `json:"mentionedPeople,omitempty"`
		Before          string `json:"before,omitempty"`
		BeforeMessage   string `json:"beforeMessage,omitempty"`
		Max             int    `json:"max,omitempty"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.RoomId == "" {
		return nil, fmt.Errorf("roomId is required")
	}

	queryParams := map[string]string{"roomId": params.RoomId}
	if params.ParentId != "" {
		queryParams["parentId"] = params.ParentId
	}
	if params.MentionedPeople != "" {
		queryParams["mentionedPeople"] = params.MentionedPeople
	}
	if params.Before != "" {
		queryParams["before"] = params.Before
	}
	if params.BeforeMessage != "" {
		queryParams["beforeMessage"] = params.BeforeMessage
	}
	if params.Max > 0 {
		queryParams["max"] = strconv.Itoa(params.Max)
	}

	return t.client.Get("/messages", queryParams)
}

func (t *ListMessagesTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// CreateMessageTool creates a new message in a Webex room
type CreateMessageTool struct {
	ToolBase
}

func NewCreateMessageTool() *CreateMessageTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId":       StringProperty("The ID of the room to send the message to."),
		"parentId":     StringProperty("The ID of the parent message (for threaded messages)."),
		"toPersonId":   StringProperty("The ID of the person to send a direct message to."),
		"toPersonEmail": StringProperty("The email address of the person to send a direct message to."),
		"text":         StringProperty("The plain text content of the message."),
		"markdown":     StringProperty("The markdown content of the message."),
		"html":         StringProperty("The HTML content of the message."),
		"files":        ArrayProperty("Array of file URLs to attach to the message.", StringProperty("")),
		"attachments":  ArrayProperty("Array of attachment objects for cards.", ObjectProperty("")),
	}, []string{})

	return &CreateMessageTool{
		ToolBase: NewToolBase("create_message", "Create a new message in a Webex room.", schema),
	}
}

func (t *CreateMessageTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Validate that at least one recipient is specified
	if params["roomId"] == nil && params["toPersonId"] == nil && params["toPersonEmail"] == nil {
		return nil, fmt.Errorf("one of roomId, toPersonId, or toPersonEmail is required")
	}

	// Validate that at least one content type is specified
	if params["text"] == nil && params["markdown"] == nil && params["html"] == nil &&
		params["files"] == nil && params["attachments"] == nil {
		return nil, fmt.Errorf("at least one of text, markdown, html, files, or attachments is required")
	}

	return t.client.Post("/messages", params)
}

func (t *CreateMessageTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// DeleteMessageTool deletes a message from Webex
type DeleteMessageTool struct {
	ToolBase
}

func NewDeleteMessageTool() *DeleteMessageTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"messageId": StringProperty("The ID of the message to delete."),
	}, []string{"messageId"})

	return &DeleteMessageTool{
		ToolBase: NewToolBase("delete_message", "Delete a message from Webex.", schema),
	}
}

func (t *DeleteMessageTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		MessageId string `json:"messageId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.MessageId == "" {
		return nil, fmt.Errorf("messageId is required")
	}

	err := t.client.Delete(fmt.Sprintf("/messages/%s", params.MessageId))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"success": true}, nil
}

func (t *DeleteMessageTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// EditMessageTool edits an existing message in Webex
type EditMessageTool struct {
	ToolBase
}

func NewEditMessageTool() *EditMessageTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"messageId": StringProperty("The ID of the message to edit."),
		"roomId":    StringProperty("The ID of the room containing the message."),
		"text":      StringProperty("The new plain text content of the message."),
		"markdown":  StringProperty("The new markdown content of the message."),
		"html":      StringProperty("The new HTML content of the message."),
	}, []string{"messageId", "roomId"})

	return &EditMessageTool{
		ToolBase: NewToolBase("edit_message", "Edit an existing message in Webex.", schema),
	}
}

func (t *EditMessageTool) Execute(args json.RawMessage) (interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	messageId, ok := params["messageId"].(string)
	if !ok || messageId == "" {
		return nil, fmt.Errorf("messageId is required")
	}

	if params["roomId"] == nil || params["roomId"] == "" {
		return nil, fmt.Errorf("roomId is required")
	}

	// At least one content type must be specified
	if params["text"] == nil && params["markdown"] == nil && params["html"] == nil {
		return nil, fmt.Errorf("at least one of text, markdown, or html is required")
	}

	// Remove messageId from params as it's in the URL
	delete(params, "messageId")

	return t.client.Put(fmt.Sprintf("/messages/%s", messageId), params)
}

func (t *EditMessageTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// GetMessageDetailsTool gets details of a specific message
type GetMessageDetailsTool struct {
	ToolBase
}

func NewGetMessageDetailsTool() *GetMessageDetailsTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"messageId": StringProperty("The ID of the message to get details for."),
	}, []string{"messageId"})

	return &GetMessageDetailsTool{
		ToolBase: NewToolBase("get_message_details", "Get details of a specific message.", schema),
	}
}

func (t *GetMessageDetailsTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		MessageId string `json:"messageId"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.MessageId == "" {
		return nil, fmt.Errorf("messageId is required")
	}

	return t.client.Get(fmt.Sprintf("/messages/%s", params.MessageId), nil)
}

func (t *GetMessageDetailsTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}

// ListDirectMessagesTool lists direct messages between two people
type ListDirectMessagesTool struct {
	ToolBase
}

func NewListDirectMessagesTool() *ListDirectMessagesTool {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"personId":      StringProperty("The ID of the person to list direct messages with."),
		"personEmail":   StringProperty("The email of the person to list direct messages with."),
		"max":           IntegerProperty("Limit the maximum number of messages in the response."),
		"before":        StringProperty("List messages sent before a specific date and time."),
		"beforeMessage": StringProperty("List messages sent before a specific message, by ID."),
	}, []string{})

	return &ListDirectMessagesTool{
		ToolBase: NewToolBase("list_direct_messages", "List direct messages with a person.", schema),
	}
}

func (t *ListDirectMessagesTool) Execute(args json.RawMessage) (interface{}, error) {
	var params struct {
		PersonId      string `json:"personId,omitempty"`
		PersonEmail   string `json:"personEmail,omitempty"`
		Max           int    `json:"max,omitempty"`
		Before        string `json:"before,omitempty"`
		BeforeMessage string `json:"beforeMessage,omitempty"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if params.PersonId == "" && params.PersonEmail == "" {
		return nil, fmt.Errorf("either personId or personEmail is required")
	}

	queryParams := map[string]string{}
	if params.PersonId != "" {
		queryParams["personId"] = params.PersonId
	}
	if params.PersonEmail != "" {
		queryParams["personEmail"] = params.PersonEmail
	}
	if params.Max > 0 {
		queryParams["max"] = strconv.Itoa(params.Max)
	}
	if params.Before != "" {
		queryParams["before"] = params.Before
	}
	if params.BeforeMessage != "" {
		queryParams["beforeMessage"] = params.BeforeMessage
	}

	return t.client.Get("/messages/direct", queryParams)
}

func (t *ListDirectMessagesTool) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}
	return t.Execute(argsJSON)
}