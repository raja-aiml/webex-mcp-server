package tools

import (
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/raja-aiml/webex-mcp-server-go/webex"
)

// Example of using generic tools for messages

// ListMessagesParams for the generic list messages tool
type ListMessagesParams struct {
	RoomId          string `json:"roomId" required:"true"`
	ParentId        string `json:"parentId,omitempty"`
	MentionedPeople string `json:"mentionedPeople,omitempty"`
	Before          string `json:"before,omitempty"`
	BeforeMessage   string `json:"beforeMessage,omitempty"`
	Max             int    `json:"max,omitempty"`
}

// NewListMessagesToolGeneric creates a list messages tool using generics
func NewListMessagesToolGeneric() Tool {
	properties := map[string]*jsonschema.Schema{
		"roomId":          StringProperty("The ID of the room to list messages from."),
		"parentId":        StringProperty("The ID of the parent message to filter by."),
		"mentionedPeople": StringProperty("List messages with these people mentioned, by ID."),
		"before":          StringProperty("List messages sent before a specific date and time."),
		"beforeMessage":   StringProperty("List messages sent before a specific message, by ID."),
		"max":             IntegerProperty("Limit the maximum number of messages in the response."),
	}

	return NewListTool[ListMessagesParams](
		"list_messages_generic",
		"List messages in a Webex room (generic implementation).",
		"/messages",
		properties,
	)
}

// CreateMessageParams for the generic create message tool
type CreateMessageParams struct {
	RoomId        string                   `json:"roomId,omitempty"`
	ParentId      string                   `json:"parentId,omitempty"`
	ToPersonId    string                   `json:"toPersonId,omitempty"`
	ToPersonEmail string                   `json:"toPersonEmail,omitempty"`
	Text          string                   `json:"text,omitempty"`
	Markdown      string                   `json:"markdown,omitempty"`
	HTML          string                   `json:"html,omitempty"`
	Files         []string                 `json:"files,omitempty"`
	Attachments   []map[string]interface{} `json:"attachments,omitempty"`
}

// NewCreateMessageToolGeneric creates a create message tool using generics
func NewCreateMessageToolGeneric() Tool {
	properties := map[string]*jsonschema.Schema{
		"roomId":        StringProperty("The ID of the room to send the message to."),
		"parentId":      StringProperty("The ID of the parent message (for threaded messages)."),
		"toPersonId":    StringProperty("The ID of the person to send a direct message to."),
		"toPersonEmail": StringProperty("The email address of the person to send a direct message to."),
		"text":          StringProperty("The plain text content of the message."),
		"markdown":      StringProperty("The markdown content of the message."),
		"html":          StringProperty("The HTML content of the message."),
		"files":         ArrayProperty("Array of file URLs to attach to the message.", StringProperty("")),
		"attachments":   ArrayProperty("Array of attachment objects for cards.", ObjectProperty("")),
	}

	schema := SimpleSchema(properties, []string{})

	// Custom executor for validation
	executor := func(params *CreateMessageParams, client webex.HTTPClient) (interface{}, error) {
		// Validate that at least one recipient is specified
		if params.RoomId == "" && params.ToPersonId == "" && params.ToPersonEmail == "" {
			return nil, fmt.Errorf("one of roomId, toPersonId, or toPersonEmail is required")
		}

		// Validate that at least one content type is specified
		if params.Text == "" && params.Markdown == "" && params.HTML == "" &&
			len(params.Files) == 0 && len(params.Attachments) == 0 {
			return nil, fmt.Errorf("at least one of text, markdown, html, files, or attachments is required")
		}

		return client.Post("/messages", params)
	}

	return NewGenericTool(
		"create_message_generic",
		"Create a new message in a Webex room (generic implementation).",
		schema,
		executor,
	)
}

// MessageIDParams for delete/get operations
type MessageIDParams struct {
	MessageId string `json:"messageId" required:"true"`
}

// NewDeleteMessageToolGeneric creates a delete message tool using generics
func NewDeleteMessageToolGeneric() Tool {
	return NewDeleteTool(
		"delete_message_generic",
		"Delete a message from Webex (generic implementation).",
		"/messages",
		"messageId",
	)
}

// NewGetMessageDetailsToolGeneric creates a get message details tool using generics
func NewGetMessageDetailsToolGeneric() Tool {
	return NewGetTool(
		"get_message_details_generic",
		"Get details of a specific message (generic implementation).",
		"/messages",
		"messageId",
	)
}

// --- Example for Rooms ---

// RoomListParams for listing rooms
type RoomListParams struct {
	TeamId string `json:"teamId,omitempty"`
	Type   string `json:"type,omitempty"`
	SortBy string `json:"sortBy,omitempty"`
	Max    int    `json:"max,omitempty"`
}

// NewListRoomsToolGeneric creates a list rooms tool using generics
func NewListRoomsToolGeneric() Tool {
	properties := map[string]*jsonschema.Schema{
		"teamId": StringProperty("List rooms associated with a team, by ID."),
		"type":   StringProperty("List rooms by type: 'direct' or 'group'."),
		"sortBy": StringProperty("Sort results by: 'id', 'lastactivity', or 'created'."),
		"max":    IntegerProperty("Limit the maximum number of rooms in the response."),
	}

	return NewListTool[RoomListParams](
		"list_rooms_generic",
		"List Webex rooms (generic implementation).",
		"/rooms",
		properties,
	)
}

// CreateRoomParams for creating rooms
type CreateRoomParams struct {
	Title              string `json:"title" required:"true"`
	TeamId             string `json:"teamId,omitempty"`
	ClassificationId   string `json:"classificationId,omitempty"`
	IsLocked           bool   `json:"isLocked,omitempty"`
	IsPublic           bool   `json:"isPublic,omitempty"`
	Description        string `json:"description,omitempty"`
	IsAnnouncementOnly bool   `json:"isAnnouncementOnly,omitempty"`
}

// NewCreateRoomToolGeneric creates a create room tool using generics
func NewCreateRoomToolGeneric() Tool {
	properties := map[string]*jsonschema.Schema{
		"title":              StringProperty("A user-friendly name for the room."),
		"teamId":             StringProperty("The ID for the team with which this room is associated."),
		"classificationId":   StringProperty("The classification ID for the room."),
		"isLocked":           BooleanProperty("Whether the room is locked (moderator approval required)."),
		"isPublic":           BooleanProperty("Whether the room is public (allows guest users)."),
		"description":        StringProperty("The description of the room."),
		"isAnnouncementOnly": BooleanProperty("Whether the room is announcement only."),
	}

	return NewCreateTool[CreateRoomParams](
		"create_room_generic",
		"Create a new Webex room (generic implementation).",
		"/rooms",
		properties,
		[]string{"title"},
	)
}
