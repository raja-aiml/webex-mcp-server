package tools

import (
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/webex"
)

// ListMessagesParams defines the parameters for listing messages
type ListMessagesParams struct {
	RoomId          string `json:"roomId" required:"true"`
	ParentId        string `json:"parentId,omitempty" query:"parentId"`
	MentionedPeople string `json:"mentionedPeople,omitempty" query:"mentionedPeople"`
	Before          string `json:"before,omitempty" query:"before"`
	BeforeMessage   string `json:"beforeMessage,omitempty" query:"beforeMessage"`
	Max             int    `json:"max,omitempty" query:"max" includeZero:"false"`
}

// CreateMessageParams defines the parameters for creating a message
type CreateMessageParams struct {
	RoomId        string                   `json:"roomId,omitempty"`
	ToPersonId    string                   `json:"toPersonId,omitempty"`
	ToPersonEmail string                   `json:"toPersonEmail,omitempty"`
	Text          string                   `json:"text,omitempty"`
	Markdown      string                   `json:"markdown,omitempty"`
	Html          string                   `json:"html,omitempty"`
	Files         []string                 `json:"files,omitempty"`
	Attachments   []map[string]interface{} `json:"attachments,omitempty"`
	ParentId      string                   `json:"parentId,omitempty"`
}

// UpdateMessageParams defines the parameters for updating a message
type UpdateMessageParams struct {
	MessageId string `json:"messageId" required:"true"`
	RoomId    string `json:"roomId" required:"true"`
	Text      string `json:"text,omitempty"`
	Markdown  string `json:"markdown,omitempty"`
}

// NewListMessagesTool lists messages in a room
func NewListMessagesTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"roomId":          StringProperty("List messages in a room, by ID."),
		"parentId":        StringProperty("List messages with a parent, by ID."),
		"mentionedPeople": StringProperty("List messages with these people mentioned."),
		"before":          StringProperty("List messages sent before a date and time (ISO8601 format)."),
		"beforeMessage":   StringProperty("List messages sent before a message, by ID."),
		"max":             IntegerProperty("Limit the maximum number of messages in the response."),
	}

	return NewListTool[ListMessagesParams](
		"list_messages",
		"List messages in a room.",
		"/messages",
		properties,
		[]string{"roomId"}, // roomId is required for listing messages
	)
}

// NewCreateMessageTool creates a new message
func NewCreateMessageTool() Tool {
	attachmentSchema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"contentType": StringProperty("The content type of the attachment."),
			"content":     ObjectProperty("The content of the attachment.", map[string]*jsonschema.Schema{}),
		},
	}

	properties := map[string]*jsonschema.Schema{
		"roomId":        StringProperty("The room ID of the message."),
		"toPersonId":    StringProperty("The person ID of the recipient when sending a 1:1 message."),
		"toPersonEmail": StringProperty("The email address of the recipient when sending a 1:1 message."),
		"text":          StringProperty("The plain text content of the message."),
		"markdown":      StringProperty("The Markdown content of the message."),
		"html":          StringProperty("The HTML content of the message."),
		"files":         ArrayProperty("File URLs to be attached to the message.", StringProperty("")),
		"attachments":   ArrayProperty("Content attachments to attach to the message.", attachmentSchema),
		"parentId":      StringProperty("The parent message to reply to."),
	}

	// Create a simple schema without oneOf (MCP doesn't support oneOf at top level)
	schema := &jsonschema.Schema{
		Type:        "object",
		Description: "Post a new message to a room or person. Specify either roomId, toPersonId, or toPersonEmail.",
		Properties:  properties,
		// No required fields at schema level - validation handled in code
	}

	return NewGenericTool("create_a_message", "Post a new message to a room or person.", schema,
		func(params *map[string]interface{}, client webex.HTTPClient) (interface{}, error) {
			// Validate that exactly one recipient field is specified
			hasRoomId := (*params)["roomId"] != nil && (*params)["roomId"] != ""
			hasToPersonId := (*params)["toPersonId"] != nil && (*params)["toPersonId"] != ""
			hasToPersonEmail := (*params)["toPersonEmail"] != nil && (*params)["toPersonEmail"] != ""

			recipientCount := 0
			if hasRoomId {
				recipientCount++
			}
			if hasToPersonId {
				recipientCount++
			}
			if hasToPersonEmail {
				recipientCount++
			}

			if recipientCount == 0 {
				return nil, fmt.Errorf("exactly one of roomId, toPersonId, or toPersonEmail is required")
			}
			if recipientCount > 1 {
				return nil, fmt.Errorf("only one of roomId, toPersonId, or toPersonEmail should be specified")
			}

			// Validate that at least one content field is specified
			hasText := (*params)["text"] != nil && (*params)["text"] != ""
			hasMarkdown := (*params)["markdown"] != nil && (*params)["markdown"] != ""
			hasHtml := (*params)["html"] != nil && (*params)["html"] != ""
			hasFiles := (*params)["files"] != nil
			hasAttachments := (*params)["attachments"] != nil

			if !hasText && !hasMarkdown && !hasHtml && !hasFiles && !hasAttachments {
				return nil, fmt.Errorf("at least one of text, markdown, html, files, or attachments is required")
			}

			return client.Post("/messages", *params)
		})
}

// NewGetMessageDetailsTool gets details of a specific message
func NewGetMessageDetailsTool() Tool {
	return NewGetTool(
		"get_message_details",
		"Get details of a message by ID.",
		"/messages",
		"messageId",
		"The unique identifier for the message.",
	)
}

// NewUpdateMessageTool updates a message
func NewUpdateMessageTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"messageId": StringProperty("The unique identifier for the message."),
		"roomId":    StringProperty("The room ID of the message."),
		"text":      StringProperty("The plain text content of the message."),
		"markdown":  StringProperty("The Markdown content of the message."),
	}

	return NewUpdateTool[UpdateMessageParams](
		"update_a_message",
		"Update a message.",
		"/messages",
		"messageId",
		properties,
		[]string{"messageId", "roomId"},
	)
}

// NewDeleteMessageTool deletes a message
func NewDeleteMessageTool() Tool {
	return NewDeleteTool(
		"delete_a_message",
		"Delete a message.",
		"/messages",
		"messageId",
		"The unique identifier for the message.",
	)
}

// NewListDirectMessagesTool lists direct messages
func NewListDirectMessagesTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"personId":    StringProperty("List messages in a 1:1 room with this person."),
		"personEmail": StringProperty("List messages in a 1:1 room with this person email."),
		"max":         IntegerProperty("Limit the maximum number of messages in the response."),
	}

	schema := SimpleSchema("List messages in a 1:1 space.", properties, []string{})

	return NewGenericTool("list_direct_messages", "List messages in a 1:1 space.", schema,
		func(params *map[string]interface{}, client webex.HTTPClient) (interface{}, error) {
			queryParams := make(map[string]string)

			if personId, ok := (*params)["personId"].(string); ok && personId != "" {
				queryParams["personId"] = personId
			}
			if personEmail, ok := (*params)["personEmail"].(string); ok && personEmail != "" {
				queryParams["personEmail"] = personEmail
			}
			if max, ok := (*params)["max"].(float64); ok && max > 0 {
				queryParams["max"] = fmt.Sprintf("%d", int(max))
			}

			return client.Get("/messages/direct", queryParams)
		})
}
