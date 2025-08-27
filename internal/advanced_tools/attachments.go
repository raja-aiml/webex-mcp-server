package advanced_tools

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/raja-aiml/webex-mcp-server/internal/tools"
)

type CreateAttachmentActionParams struct {
	Type      string                 `json:"type" required:"true"`
	MessageId string                 `json:"messageId" required:"true"`
	Inputs    map[string]interface{} `json:"inputs,omitempty"`
}

// NewCreateAttachmentActionTool creates an attachment action
func NewCreateAttachmentActionTool() Tool {
	properties := map[string]*jsonschema.Schema{
		"type":      StringProperty("The type of action."),
		"messageId": StringProperty("The ID of the message with attachment."),
		"inputs":    ObjectProperty("The attachment action's inputs.", map[string]*jsonschema.Schema{}),
	}

	return tools.NewCreateTool[CreateAttachmentActionParams](
		"create_an_attachment_action",
		"Create an attachment action",
		"/attachment/actions",
		properties,
		[]string{"type", "messageId"},
	)
}

// NewGetAttachmentActionDetailsTool gets attachment action details
func NewGetAttachmentActionDetailsTool() Tool {
	return tools.NewGetTool(
		"get_attachment_action_details",
		"Get details for an attachment action by ID.",
		"/attachment/actions",
		"attachmentActionId",
		"The unique identifier for the attachment action.",
	)
}
