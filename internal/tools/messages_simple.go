package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// Example of simplified message tools using ToolBase to demonstrate DRY principle

// ListMessagesSimple - simplified implementation
type ListMessagesSimple struct {
	ToolBase
}

// NewListMessagesSimple creates a simplified list messages tool
func NewListMessagesSimple() *ListMessagesSimple {
	schema := SimpleSchema(map[string]*jsonschema.Schema{
		"roomId":          StringProperty("The ID of the room to list messages from."),
		"parentId":        StringProperty("The ID of the parent message to filter by."),
		"mentionedPeople": StringProperty("List messages with these people mentioned, by ID."),
		"before":          StringProperty("List messages sent before a specific date and time."),
		"beforeMessage":   StringProperty("List messages sent before a specific message, by ID."),
		"max":             IntegerProperty("Limit the maximum number of messages in the response."),
	}, []string{"roomId"})

	return &ListMessagesSimple{
		ToolBase: NewToolBase("list_messages", "List messages in a Webex room.", schema),
	}
}

// Execute implements the tool execution
func (t *ListMessagesSimple) Execute(args json.RawMessage) (interface{}, error) {
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

	// Build query parameters
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
	} else {
		queryParams["max"] = "100"
	}

	return t.client.Get("/messages", queryParams)
}

// ExecuteWithMap delegates to Execute
func (t *ListMessagesSimple) ExecuteWithMap(args map[string]interface{}) (interface{}, error) {
	return ExecuteWithMapBase(t, args)
}
